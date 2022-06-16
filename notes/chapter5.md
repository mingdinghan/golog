## Chapter 5 - Secure Your Services

- Security in Distributed Services
  1. Encrypt data in-flight to protect against man-in-the-middle attacks -> TLS (include server authentication)
  2. Authentication to identify clients (mutual TLS is commonly used for machine-to-machine communications)
  3. Authorization to determine the permissions of the identified clients (ACL: access-control lists)
- Operate as Your Own Certificate Authority with CFSSL
  - good for internal services, instead of going through a third-party CA
  - `cfssl` is used to sign, verify, and bundle TLS certificates and output the results as JSON
  - `cfssljson` is used to take the JSON output and split them into separate key, certificate, CSR, and bundle files
    ```bash
    $ go get github.com/cloudflare/cfssl/cmd/cfssl@v1.6.1
    $ go get github.com/cloudflare/cfssl/cmd/cfssljson@v1.6.1
    ```
- define variables to specify the paths to the generated TLS certs in order to look up and parse for tests
  - use the cert and key files to build `*tls.Configs`
- Access Control List: a table of rules where each row says "Subject A is permitted to do action B on object C"
- Casbin is a library that supports enforcing authorization based on various control models, like ACL and RBAC
  ```bash
  $ go get github.com/casbin/casbin/v2
  ```
  - the `model` configures Casbin to use ACL as its authorization mechanism
  - the `policy` is a CSV file containing the ACL table
