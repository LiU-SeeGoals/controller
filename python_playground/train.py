import torch
import torch.optim as optim
from models import SeeGoalsDNN


def goto_loss(delta, positions, target = torch.tensor([0.0, 0.0])):
    """Minimize the distance between the predicted position and the target position"""
    predicted_positions = torch.abs(positions + delta)
    distance = torch.norm(predicted_positions - target)
    original_distance = torch.norm(positions - target)
    return torch.mean(distance - original_distance)

def apply_force_bese_on_distance(ball, team):
    """Apply a force to the team based on the distance to the ball"""
    vec =  ball - team
    dist = torch.norm(vec, dim=1) 
    scaling = (1/(dist**(1/5))).unsqueeze(1)
    # print(f"{scaling=}")
    force = torch.sum(vec * scaling, dim=1)
    # print(f"{force=}")
    new_pos_ball = ball + force
    # print(f"{team=}, {ball=}, {new_pos_ball=}, {dist=}, {force=} ")
    return new_pos_ball


def ball_to_target_loss(my_team, delta, ball, target = torch.tensor([0.0, 0.0])):
    """Minimize the distance between the predicted position and the target position"""
    scale = 100
    ball_pos = apply_force_bese_on_distance(ball, my_team)
    original_distance = torch.norm(ball_pos - target)

    predicted_positions = my_team + delta * scale
    ball_pos_predicted = apply_force_bese_on_distance(ball, predicted_positions)
    distance = torch.norm(ball_pos_predicted - target)
    # print(distance, original_distance, distance - original_distance)
    # 1/0
    return torch.mean(distance - original_distance)


    

# Instantiate model, optimizer, and loss function
field_hight = 9000
field_width = 7000
players_per_team = 2
batch_size = 100
model = SeeGoalsDNN(num_players_per_team=players_per_team, field_hight=field_hight, field_width=field_width)
optimizer = optim.Adam(model.parameters(), lr=0.00001)


# Training loop
epochs = 5000
losses = []
for epoch in range(epochs):
    my_team = torch.rand(batch_size, players_per_team, 6)-0.5
    my_team[:, :, :2] *= torch.tensor([field_width, field_hight])
    enemy_team = torch.rand(batch_size, players_per_team, 6)-0.5
    enemy_team[:, :, :2] *= torch.tensor([field_width, field_hight])
    ball = torch.rand(batch_size,1, 6)
    ball[:, :, :2] *= torch.tensor([field_width, field_hight]) * 0.2
    # Set gradients to zero
    optimizer.zero_grad()

    # Forward pass
    predicted_delta = model(my_team, enemy_team, ball)

    # Loss calculation
    loss = goto_loss(predicted_delta, my_team[:, :, :2])
    # loss = ball_to_target_loss(my_team[:, :, :2], predicted_delta[:, :, :2], ball[:, :, :2])
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