# Maintainer: Your Name <your_email@example.com>
BUILDDIR=/tmp/makepkg-git-wrap

pkgname=git-wrap
pkgver=0.0.1
pkgrel=1
pkgdesc="An EC-compliant Git and Submodule manager CLI wrapper"
arch=('x86_64')
url="https://github.com/Nexus29/git-wrap"
license=('MIT')
depends=('git')
provides=('git-wrap')
conflicts=('git-wrap')

package() {
    # Grabs the compiled binary from GoReleaser and copies it safely
    install -Dm755 "${startdir}/dist/git-wrap_linux_amd64_v1/git-wrap" "${pkgdir}/usr/bin/git-wrap"
}
