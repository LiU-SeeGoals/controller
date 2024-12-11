import torch
import torch.nn as nn
import torch.nn.functional as F
import numpy as np
import random

class Net(nn.Module):

    def __init__(self):
        super(Net, self).__init__()

        self.fc1 = nn.Linear(24, 72)
        self.fc2 = nn.Linear(72, 72)
        self.fc3 = nn.Linear(72, 48)
        self.fc4 = nn.Linear(48, 48)
        self.fc5 = nn.Linear(48, 4)

    def forward(self, input):
        
       
        x = torch.reshape(input, (-1, 24))
        
        f1 = F.relu(self.fc1(x))
        
        f2 = F.relu(self.fc2(f1))

        f3 = F.relu(self.fc3(f2))

        f4 = F.relu(self.fc4(f3))
          
        output = self.fc5(f4)
        
        return output

class Dataset(torch.utils.data.Dataset):
  'Characterizes a dataset for PyTorch'
  def __init__(self, X, labels):
        'Initialization'
        self.labels = labels
        self.X = X

  def __len__(self):
        'Denotes the total number of samples'
        return len(self.X)

  def __getitem__(self, index):
        'Generates one sample of data'
        # Select sample
        X = self.X[index]

        # Load data and get label
        
        y = self.labels[index]

        return X, y

def create_mock_data(x_lim, y_lim):
    all_pos = []
    y_vals = []

    for i in range(10000):
        pos = []
        count1, count2, count3, count4 = 0,0,0,0
        for j in range(12):
            x = np.random.randint(0, x_lim)
            y = np.random.randint(0, y_lim)
            pos.append([x,y])
            if x < (x_lim/2) and y < (y_lim/2):
                count1 += 1
            if x > (x_lim/2) and y < (y_lim/2):
                count2 += 1
            if x < (x_lim/2) and y > (y_lim/2):
                count3 += 1
            if x > (x_lim/2) and y > (y_lim/2):
                count4 += 1
        all_pos.append(pos)

        y_vals.append(np.array([count1, count2, count3, count4]).argmin())

    return all_pos, y_vals

def create_data_loaders(all_pos, y_vals):
    # Parameters
    params = {'batch_size': 64}


    # Datasets
    partition_train = torch.tensor(all_pos[0:8000], dtype=torch.float32)# IDs
    labels_train = torch.tensor(y_vals[0:8000])# Labels

    partition_val = torch.tensor(all_pos[8000:], dtype=torch.float32)# IDs
    labels_val = torch.tensor(y_vals[8000:])# Labels

    # Generators
    training_set = Dataset(partition_train, labels_train)
    training_generator = torch.utils.data.DataLoader(training_set, **params)

    test_set = Dataset(partition_val, labels_val)
    testing_generator = torch.utils.data.DataLoader(test_set, **params)    
    
    return training_generator, testing_generator



def train(model, training_generator, criterion, optimizer):
    model.train()

    for data in training_generator:  # Iterate in batches over the training dataset.
        x = data[0]
        y = data[1]
        print(x)
        print(y)
        out = model(x)  
        loss = criterion(out, y)  # Compute the loss.
        loss.backward()  # Derive gradients.
        optimizer.step()  # Update parameters based on gradients.
        optimizer.zero_grad()  # Clear gradients.

def test(model, loader):
    model.eval()

    correct = 0
    for data in loader:  # Iterate in batches over the training/test dataset.
        x = data[0]
        y = data[1]
        out = model(x)
        pred = out.argmax(dim=1)  # Use the class with highest probability.
        correct += int((pred == y).sum())  # Check against ground-truth labels.

    return correct / len(loader.dataset)  # Derive ratio of correct predictions.


def train_model(training_generator, testing_generator):
   
    model = Net()
    optimizer = torch.optim.Adam(model.parameters(), lr=0.01)
    criterion = torch.nn.CrossEntropyLoss()
    for epoch in range(1,300):
        train(model, training_generator, criterion, optimizer)
        train_acc = test(model, training_generator)
        test_acc = test(model, testing_generator)
        print(f'Epoch: {epoch:03d}, Train Acc: {train_acc:.4f}, Test Acc: {test_acc:.4f}')
    torch.save(model.state_dict(), "model.pt")


def get_model(name):
    model = Net()
    model.load_state_dict(torch.load(name, weights_only=True))
    model.eval()
    return model

def predict(model, locations):
    locations = torch.tensor(locations, dtype=torch.float32)
    out = model(locations)
    pred = out.argmax(dim=1)
    print(out)
    return pred

if __name__ == '__main__':
    input = input("train (1) or load (2)?")
    if int(input) == 1:
        all_pos, y_vals = create_mock_data(90,120)
        training_generator, testing_generator = create_data_loaders(all_pos, y_vals)
        train_model(training_generator, testing_generator)
       

    if int(input) == 2: 
        model = Net()
        model.load_state_dict(torch.load("model.pt", weights_only=True))
        model.eval()   
        out = model(torch.tensor([[1,2], [2,3], [78,90], [90, 1], [32,45], [84,32], [1,1],[12,109], [12,34], [90,90], [34, 52], [10, 100]], dtype=torch.float32))
        print(out)  
        pred = out.argmax(dim=1)
        print(pred)
    