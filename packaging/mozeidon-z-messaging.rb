class MozeidonZMessaging < Formula
  desc "Native-messaging host bridging the Mozeidon-Z browser extension and CLI"
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

  def caveats
    <<~EOS
      Firefox needs a native-messaging host manifest to launch this bridge.
      Homebrew can't write outside its prefix, so install it once per macOS
      user (the Mozeidon-Z toolbar popup shows "Disconnected" until you do):

        mkdir -p ~/Library/Application\\ Support/Mozilla/NativeMessagingHosts
        cat > ~/Library/Application\\ Support/Mozilla/NativeMessagingHosts/mozeidon_z.json <<'JSON'
        {"name":"mozeidon_z","description":"Mozeidon-Z native messaging host","path":"#{opt_bin}/mozeidon-z-messaging","type":"stdio","allowed_extensions":["mozeidon-z@a-layer.io"]}
        JSON

      Then install the Mozeidon-Z add-on from addons.mozilla.org. No Firefox
      restart is needed — the extension reconnects within ~1s.
    EOS
  end

  test do
    assert_match "mozeidon-z-messaging", shell_output("#{bin}/mozeidon-z-messaging --version")
  end
end
