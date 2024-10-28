import torch
import torch.optim as optim
from models import SeeGoalsDNN



def apply_force_based_on_distance(ball, team):
    """Apply a force to the ball based on the distance to the team."""
    # Compute vectors from team players to the ball
    vec = team - ball
    dist = torch.norm(vec, dim=2)  
    # Compute scaling factor based on the distance
    scaling = (1 / (dist ** (1/2))).unsqueeze(2)
    force = torch.sum(vec * scaling, dim=1).unsqueeze(1)
    # Update ball positions by adding the computed force
    new_pos_ball = ball + force  
    return new_pos_ball

def ball_to_target_loss(my_team, delta, ball, target = torch.tensor([0.0, 0.0])):
    """Minimize the distance between the ball's predicted position and the target position."""
    scale = 1
    batch_size = my_team.shape[0]
    num_players = my_team.shape[1]
    # ball_pos = apply_force_based_on_distance(ball, my_team)  
    # target = target.repeat(batch_size, 1).unsqueeze(1)
    # ball_to_target = torch.norm(ball - target, dim=2)
    predicted_positions = my_team + delta * scale  
    original_team_to_ball = torch.mean(torch.norm(my_team - ball, dim=2), dim=1)
    predicted_team_to_ball = torch.mean(torch.norm(predicted_positions - ball, dim=2), dim=1)
    # ball_pos_predicted = apply_force_based_on_distance(ball, predicted_positions)  
    # pred_dist_to_target = torch.norm(ball_pos_predicted - target, dim=1) 

    # Compute the loss as the mean difference in distances
    # print(pred_dist_to_target.shape, ball_to_target.shape, team_to_ball.shape, target.shape, ball_pos_predicted.shape)
    # loss = torch.mean(pred_dist_to_target - ball_to_target) + torch.mean(team_to_ball)
    team_to_ball_loss = torch.mean(predicted_team_to_ball - original_team_to_ball)
    return team_to_ball_loss

    
def train(model = None, 
        field_hight = 9000,
        field_width = 7000,
        players_per_team = 2,
        batch_size = 256,
        lr = 0.0001,
        epochs = 20_000,
        loss_fn = ball_to_target_loss,
        load_model = False
        ):
    scale_field = torch.tensor([field_width, field_hight]).float()
    if model is None:
        model = SeeGoalsDNN(num_players_per_team=players_per_team, field_hight=field_hight, field_width=field_width)
    if load_model:
        model.load_state_dict(torch.load(model.path))
        model.eval()
        return model
    optimizer = optim.Adam(model.parameters(), lr=lr)
    model.train()


    # Training loop
    losses = []
    for epoch in range(epochs):
        my_team = torch.rand(batch_size, players_per_team, 6)-0.5
        my_team[:, :, :2] *= scale_field
        enemy_team = torch.rand(batch_size, players_per_team, 6)-0.5
        enemy_team[:, :, :2] *= scale_field
        ball = torch.rand(batch_size, 1, 6)-0.5
        ball[:, :, :2] *= scale_field
        # Set gradients to zero
        optimizer.zero_grad()

        # Forward pass
        predicted_delta = model(my_team, enemy_team, ball)

        # Loss calculation
        loss = loss_fn(my_team[:, :, :2], predicted_delta[:, :, :2], ball[:, :, :2])
        # Backward pass and optimization
        loss.backward()
        optimizer.step()
        losses.append(loss.item())

        # Print the loss for every 50 epochs
        if epoch % 50 == 0:
            mean = sum(losses[-50:]) / 50
            print(f"Epoch [{epoch}/{epochs}], Loss: {mean:.5f}")
            torch.save(model.state_dict(), model.path)
    model.eval()
    return model

if __name__ == '__main__':
    train()