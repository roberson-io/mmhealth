# Start from the Debian base image
FROM debian:trixie-slim

# Install necessary packages
RUN apt-get update && apt-get install -y \
    pandoc \
    perl \
    curl \
    fontconfig \
    libgcc1 \
    gnupg \
    wget \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

ARG TL_MIRROR="https://mirror.ctan.org/systems/texlive/tlnet"
ARG TL_YEAR="2025"

RUN mkdir "/tmp/texlive"
WORKDIR "/tmp/texlive"
RUN wget "$TL_MIRROR/install-tl-unx.tar.gz"

RUN tar xzvf ./install-tl-unx.tar.gz && \
    ( \
        echo "selected_scheme scheme-basic" && \
        echo "instopt_adjustpath 0" && \
        echo "tlpdbopt_install_docfiles 0" && \
        echo "tlpdbopt_install_srcfiles 0" && \
        echo "TEXDIR /opt/texlive/${TL_YEAR}" && \
        echo "TEXMFLOCAL /opt/texlive/texmf-local" && \
        echo "TEXMFSYSCONFIG /opt/texlive/${TL_YEAR}/texmf-config" && \
        echo "TEXMFSYSVAR /opt/texlive/${TL_YEAR}/texmf-var" && \
        echo "TEXMFHOME ~/.texmf" \
    ) > "/tmp/texlive.profile"

RUN cd install-tl-* && ./install-tl --location "$TL_MIRROR" -profile "/tmp/texlive.profile" --no-interaction
RUN rm -vf "/opt/texlive/${TL_YEAR}/install-tl"
RUN rm -vf "/opt/texlive/${TL_YEAR}/install-tl.log"
RUN rm -vrf /tmp/*

ENV PATH="/opt/texlive/${TL_YEAR}/bin/x86_64-linux:${PATH}"

# Update tlmgr and install tabularray
RUN tlmgr update --self && tlmgr install tabularray etoolbox xcolor \
    ninecolors fontspec luatex soul roboto xkeyval xstring pgf seqsplit

# Add Maintainer Info
LABEL maintainer="colton.shaw@mattermost.com"
LABEL org.opencontainers.image.title="mmhealth"
LABEL org.opencontainers.image.description="Mattermost Healthcheck Image"
LABEL org.opencontainers.image.authors="Colton Shaw"
LABEL org.opencontainers.image.source="https://github.com/coltoneshaw/mmhealth"
LABEL org.opencontainers.image.licenses=MIT

RUN mkdir /files

# Set the Current Working Directory inside the container
WORKDIR /app

COPY ./template/template.tex .
COPY ./template/*.png .

ENV TEXMFVAR=/tmp/texmf-var
RUN mkdir -p /tmp/texmf-var
RUN chmod -R 777 /tmp/texmf-var

# ENV TEXMFHOME /opt/texlive/${TL_YEAR}/texmf-dist

# Print the PATH environment variable
RUN echo $PATH
RUN ls /opt/texlive/${TL_YEAR}/bin/x86_64-linux


ENTRYPOINT ["pandoc", "--pdf-engine=lualatex"]