# Contributing to OpenKasugai Controller

This controller is published as a reference for proposing [DCI Concept](https://www.ntt-review.jp/archive/ntttechnical.php?contents=ntr202401fa8.html) to related existing OSS projects such as CNCF. In the future, discussions and contributions are planned within existing OSS projects.

On the other hand, if you have any specific requests for contributions to OpenKasugai Controller, please use the [Discussions](https://github.com/openkasugai/controller/discussions) feature on GitHub to make them.

## Reporting bugs and requests for additional functions

Use the [Issues](https://github.com/openkasugai/controller/issues) on GitHub.

### BUGS

- Apply the label `bug` to the issue.
- Please include the following in the comments at the top:
  - Environment
    - OpenKasugai Controller version
    - Node configuration, OS Kernel version, device model number, etc.
  - Problem Event
    - Describe the specific situation where the problem occurred.
  - Reproduction Step (optional)
    - If the reproduction procedure is complex, please describe it.

### Added Functions

- Apply the label `enhancement` to the issue.
- Please include the following in the comments at the top:
  - Functional Specifications
    - Provide final specifications to be implemented.
      - If the specification is changed by discussion, make sure that the final specification is visible in the comments at the top.
  - Reference
    - Provide links to discussions and other resources that help you understand the discussion.

### Documentation Improvements

- Apply the label `documentation` to the issue.
- Please include the following in the comments at the top:
  - Overview
    - Describe the document to be changed and the modification.

## Other questions, design requests, ideas, etc.

Use the [Discussions](https://github.com/openkasugai/controller/discussions) feature on GitHub

## Pull Requests

- At first fork the main branch of the OpenKasugai Controller repository.
- Create a topic branch on the forked repository and pull request to the main branch of the repository under OpenKasugai Controller.
  - Topic branch can have any branch name.
- You must agree to [DCO](https://developercertificate.org/) to contribute to OpenKasugai Controller.
  - Add the following signature to the commit message to indicate that you agree with the DCO.
    - `Signed-off-by: Random J Developer <random@developer.example.org>`
      - Use your real name in the signature.
      - You need to set the same name and the email in GitHub Profile.
      - `git commit -s` can add the signature.
- Associate a pull request with the corresponding Issue.
  - If there is no corresponding issue, create a new one before creating a pull request.
- Use the templates when creating a pull request.
- The title of a pull request should include "fix" followed by the issue number and a summary of the pull request.
  - `fix #[issue number] [summary]`
- Pull Request body should use a template.
- Merging Pull Requests requires a successful style check (GitHub Actions: golangci-lint) and approval from at least two maintainers.

## Release Cycle

If there are any changes such as new functions and bug fixes, we release them in October and April every year.
If there are no changes, do not release.
