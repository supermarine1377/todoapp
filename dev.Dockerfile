FROM --platform=arm64 golang:1.23

RUN useradd -ms /bin/bash dev && \
    chown -R dev:dev /go && \
    chown dev:dev /usr/local/bin

USER dev

# Install Git completion for the dev user with root privilege since this image is for local development
RUN curl -O https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh && \
    curl -O https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash && \
    chmod a+x git*.* && \
    ls -l $PWD/git*.* | awk '{print "source "$9}' >> ~/.bashrc && \
    echo "GIT_PS1_SHOWDIRTYSTATE=true" >> ~/.bashrc && \
    echo "GIT_PS1_SHOWUNTRACKEDFILES=true" >> ~/.bashrc && \
    echo "GIT_PS1_SHOWUPSTREAM=auto" >> ~/.bashrc && \
    echo 'export PS1="\[\033[01;32m\]\u@\h\[\033[01;33m\] \w \[\033[01;31m\]\$(__git_ps1 \"(%s)\") \\n\[\033[01;34m\]\\$ \[\033[00m\]"' >> ~/.bashrc

# Install development dependencies    
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install go.uber.org/mock/mockgen@v0.5.0 && \
    curl -L https://github.com/sqldef/sqldef/releases/download/v0.17.23/sqlite3def_linux_amd64.tar.gz \
        -o sqlite3def_linux_amd64.tar.gz && \
    tar xz -C /usr/local/bin -f sqlite3def_linux_amd64.tar.gz && \
    go install -v github.com/cweill/gotests/gotests@v1.6.0 && \
    go install github.com/swaggo/swag/cmd/swag@latest
  
ENV DATABASE_DSN="/app/_data/sqlite.db"
ENV TZ="Asia/Tokyo"

WORKDIR /app/
