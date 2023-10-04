#!/bin/bash

BINARY_NAME=punq-dev
VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))
VERSIONWITHOUTV=$(echo $VERSION | cut -c 2-)
SHA256_DARWIN_ARM64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-darwin-arm64.tar.gz | awk '{print $1}')
SHA256_DARWIN_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-darwin-amd64.tar.gz | awk '{print $1}')
SHA256_LINUX_ARM64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-arm64.tar.gz | awk '{print $1}')
SHA256_LINUX_ARM=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-arm.tar.gz | awk '{print $1}')
SHA256_LINUX_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-amd64.tar.gz | awk '{print $1}')
SHA256_LINUX_386=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-386.tar.gz | awk '{print $1}')
SHA256_WIN_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-windows-amd64 | awk '{print $1}')
WIN_AMD64="$BINARY_NAME-$VERSION-windows-amd64"

# Generate formula from template with replacements
cat <<EOF > punq-dev.rb
class PunqDev < Formula
  desc "View your kubernetes workloads relatively neat! [dev]"
  homepage "https://punq.dev"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/mogenius/homebrew-punq-dev/releases/download/${VERSION}/punq-dev-${VERSION}-darwin-arm64.tar.gz"
      sha256 "$SHA256_DARWIN_ARM64"
    elsif Hardware::CPU.intel?
      url "https://github.com/mogenius/homebrew-punq-dev/releases/download/${VERSION}/punq-dev-${VERSION}-darwin-amd64.tar.gz"
      sha256 "$SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/homebrew-punq-dev/releases/download/${VERSION}/punq-dev-${VERSION}-linux-amd64.tar.gz"
        sha256 "$SHA256_LINUX_AMD64"
      else
        url "https://github.com/mogenius/homebrew-punq-dev/releases/download/${VERSION}/punq-dev-${VERSION}-linux-386.tar.gz"
        sha256 "$SHA256_LINUX_386"
      end
    elif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/homebrew-punq-dev/releases/download/${VERSION}/punq-dev-${VERSION}-linux-arm64.tar.gz"
        sha256 "$SHA256_LINUX_ARM64"
      else
        url "https://github.com/mogenius/homebrew-punq-dev/releases/download/${VERSION}/punq-dev-${VERSION}-linux-arm.tar.gz"
        sha256 "$SHA256_LINUX_ARM"
      end
    end
  end
  
  version "${VERSIONWITHOUTV}"
  license "MIT"

  def install
  if OS.mac?
    if Hardware::CPU.arm?
      # Installation steps for macOS ARM64
      bin.install "punq-dev-$VERSION-darwin-arm64" => "punq-dev"
    elsif Hardware::CPU.intel?
      # Installation steps for macOS AMD64
      bin.install "punq-dev-$VERSION-darwin-amd64" => "punq-dev"
    end
  elsif OS.linux?
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux AMD64
        bin.install "punq-dev-$VERSION-linux-amd64" => "punq-dev"
      else
        # Installation steps for Linux 386
        bin.install "punq-dev-$VERSION-linux-386" => "punq-dev"
      end
    elsif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux ARM64
        bin.install "punq-dev-$VERSION-linux-arm64" => "punq-dev"
      else
        # Installation steps for Linux ARM
        bin.install "punq-dev-$VERSION-linux-arm" => "punq-dev"
      end
    end
  end
end
end
EOF

cat <<EOF > punq-dev.json
{
    "version": "$VERSIONWITHOUTV",
    "license": "MIT",
    "homepage": "https://punq.dev",
    "bin": "punq-dev.exe",
    "pre_install": "Rename-Item \"\$dir\\\\$WIN_AMD64\" punq-dev.exe",
    "description": "View your kubernetes workloads relatively neat!",
    "architecture": {
        "64bit": {
            "url": "https://github.com/mogenius/punq/releases/download/$VERSION/punq-$VERSION-windows-amd64",
            "hash": "$SHA256_WIN_AMD64"
        }
    }
}
EOF