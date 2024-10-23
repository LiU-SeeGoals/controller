import torch
import torch.optim as optim
from models import SeeGoalsDNN


def center_loss(delta, positions):
    predicted_positions = torch.abs(positions + delta)

    return torch.mean(torch.norm(predicted_positions, dim=-1))

# Instantiate model, optimizer, and loss function
field_hight = 9000
field_width = 7000
players_per_team = 2
batch_size = 100
model = SeeGoalsDNN(num_players_per_team=players_per_team, field_hight=field_hight, field_width=field_width)
optimizer = optim.Adam(model.parameters(), lr=0.001)


# Training loop
epochs = 10000
losses = []
for epoch in range(epochs):
    my_team = torch.rand(batch_size, players_per_team, 6)-0.5
    my_team[:, :, :2] *= torch.tensor([field_width, field_hight])
    enemy_team = torch.rand(batch_size, players_per_team, 6)-0.5
    enemy_team[:, :, :2] *= torch.tensor([field_width, field_hight])
    ball = torch.rand(batch_size,1, 6)
    ball[:, :, :2] *= torch.tensor([field_width, field_hight])
    # Set gradients to zero
    optimizer.zero_grad()

    # Forward pass
    predicted_delta = model(my_team, enemy_team, ball)

    # Loss calculation
    loss = center_loss(predicted_delta, my_team[:, :, :2]/torch.tensor([field_width, field_hight])*2)

    # Backward pass and optimization
    loss.backward()
    optimizer.step()
    losses.append(loss.item())

    # Print the loss for every 50 epochs
    if epoch % 50 == 0:
        mean = sum(losses[-50:]) / 50
        print(f"Epoch [{epoch}/{epochs}], Loss: {mean:.5f}")

# Save the model
torch.save(model.state_dict(), model.path)