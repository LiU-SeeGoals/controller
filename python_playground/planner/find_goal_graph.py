import torch
import torch.nn.functional as F
from torch_geometric.nn import GCNConv, MessagePassing
from torch_geometric.data import Data

class GNNEncoder(torch.nn.Module):
    def __init__(self, input_dim, hidden_dim, output_dim):
        super(GNNEncoder, self).__init__()
        self.conv1 = GCNConv(input_dim, hidden_dim)
        self.conv2 = GCNConv(hidden_dim, output_dim)

    def forward(self, x, edge_index):
        
        x = F.relu(self.conv1(x, edge_index))
        x = self.conv2(x, edge_index)
        return x

class GraphDecoder(torch.nn.Module):
    def __init__(self, hidden_dim):
        super(GraphDecoder, self).__init__()
        self.fc = torch.nn.Linear(hidden_dim, 1)  # Predict edge weights or existence

    def forward(self, x):
        # Compute pairwise edge probabilities
        return torch.mm(x, x.T)

class GraphToGraph(torch.nn.Module):
    def __init__(self, input_dim, hidden_dim, output_dim):
        super(GraphToGraph, self).__init__()
        self.encoder = GNNEncoder(input_dim, hidden_dim, output_dim)
        self.decoder = GraphDecoder(output_dim)

    def forward(self, x, edge_index):
        latent = self.encoder(x, edge_index)  # Encode graph
        output = self.decoder(latent)  # Decode to new graph structure
        return output


def compute_loss(pred_graph, target_graph):
    """
    Compute loss for graph regression.

    Args:
        pred_graph (torch.Tensor): Predicted graph-level outputs.
        target_graph (torch.Tensor): Ground-truth graph-level values.

    Returns:
        torch.Tensor: Loss value.
    """
    loss_fn = torch.nn.MSELoss()
    loss = loss_fn(pred_graph, target_graph)
    return loss

# Example usage:
# input_dim = number of features per node
# Create input graph with node features and edge list
num_nodes = 12
input_dim = 4

x = torch.rand((num_nodes, input_dim))  # Node features

target = torch.rand(num_nodes**2) 
target = torch.reshape(target, (num_nodes, num_nodes))
edges = []

for i in range(num_nodes):
    for j in range(num_nodes):
        if i != j:
            edges.append((i,j))
        

edge_index = torch.tensor(edges)  # Edges
edge_index = torch.reshape(edge_index, (2,-1))

model = GraphToGraph(input_dim, hidden_dim=64, output_dim=32)
optimizer = torch.optim.Adam(model.parameters(), lr=0.01)


params = {'batch_size': 64}
training_set = Dataset(partition_train, labels_train)
training_generator = torch.utils.data.DataLoader(training_set, **params)

test_set = Dataset(partition_val, labels_val)
testing_generator = torch.utils.data.DataLoader(test_set, **params)

epochs = 300

# Training loop
for epoch in range(epochs):
    optimizer.zero_grad()
    pred = model(x, edge_index)

    loss = compute_loss(pred, target)  # Define target graph's adjacency matrix or features
    print(loss)
    loss.backward()
    optimizer.step()
