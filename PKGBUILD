# Maintainer: Your Name <your_email@example.com>
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
    # 💡 This grabs the exact binary GoReleaser just built
    # and tells pacman to install it globally into /usr/bin/
    install -Dm755 "${srcdir}/../dist/git-wrap_linux_amd64_v1/git-wrap" "${pkgdir}/usr/bin/git-wrap"
}
