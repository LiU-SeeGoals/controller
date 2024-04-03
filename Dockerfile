FROM golang:1.21-alpine

WORKDIR /var/controller

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Install zsh and git (git might be needed for oh-my-zsh and its plugins)
RUN apk add --no-cache zsh git curl

# Install oh-my-zsh
RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended

# Install zsh-syntax-highlighting plugin
RUN git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

# Add zsh-syntax-highlighting plugin to your .zshrc
RUN sed -i 's/plugins=(git)/plugins=(git zsh-syntax-highlighting)/' ~/.zshrc

# Set zsh as default shell
CMD ["zsh"]