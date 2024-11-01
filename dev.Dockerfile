FROM --platform=arm64 golang:1.23

# Install Git complement
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
    go install github.com/golang/mock/mockgen@v1.6.0 && \
    curl -L https://github.com/sqldef/sqldef/releases/download/v0.17.23/sqlite3def_linux_amd64.tar.gz | tar xz -C /usr/local/bin

WORKDIR /app/
