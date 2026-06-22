class MozeidonZMessaging < Formula
  desc "Mozeidon-Z native-messaging host — browser ⇄ CLI IPC bridge"
  homepage "https://github.com/colangelo/mozeidon-z-messaging"
  version "VERSION_PLACEHOLDER"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/colangelo/mozeidon-z-messaging/releases/download/v#{version}/mozeidon-z-messaging-darwin-arm64"
      sha256 "DARWIN_ARM64_PLACEHOLDER"
    else
      url "https://github.com/colangelo/mozeidon-z-messaging/releases/download/v#{version}/mozeidon-z-messaging-darwin-amd64"
      sha256 "DARWIN_AMD64_PLACEHOLDER"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/colangelo/mozeidon-z-messaging/releases/download/v#{version}/mozeidon-z-messaging-linux-arm64"
      sha256 "LINUX_ARM64_PLACEHOLDER"
    else
      url "https://github.com/colangelo/mozeidon-z-messaging/releases/download/v#{version}/mozeidon-z-messaging-linux-amd64"
      sha256 "LINUX_AMD64_PLACEHOLDER"
    end
  end

  def install
    binary = Dir["mozeidon-z-messaging-*"].first || "mozeidon-z-messaging"
    bin.install binary => "mozeidon-z-messaging"
  end

  test do
    assert_match "mozeidon-z-messaging", shell_output("#{bin}/mozeidon-z-messaging --version")
  end
end
