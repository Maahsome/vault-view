# RELEASE Process

1. Create a new branch for release.
2. export RELEASE_VERSION=<newSemVer>
3. make changelog
4. Review changelogs/releases/<newSemVer>.md
5. Create and Merge PR
6. Pull main
7. git tag <newSemVer>
8. git push origin <newSemVer>
9. Check to see if Release GitHub Action kicks off
10. Check Releases Page, and homebrew-tap/Formula/vault-view.rb
