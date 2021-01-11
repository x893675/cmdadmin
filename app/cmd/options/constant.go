package options

const (
	// CertificatesDir flag sets the path where to save and read the certificates.
	CertificatesDir = "cert-dir"

	// CfgPath flag sets the path to cmdadmin config file.
	CfgPath = "config"

	// CSROnly flag instructs cmdadmin to create CSRs instead of automatically creating or renewing certs
	CSROnly = "csr-only"

	// CSRDir flag sets the location for CSRs and flags to be output
	CSRDir = "csr-dir"
)

const (
	// certs default config file
	DefaultCertsConfig = `
# # name 是该证书的名字
# # caName 是签发该证书的证书名, 如果该证书是非 CA 证书，则值置空
# # example:
# #         name: ca
# #         caName: ""
# - name: example
#   caName: "example-ca"
# # baseName 是生成证书的相应前缀, 下面的配置会在当前目录生成的 example 文件夹下生成
# # server.crt 和 server.key 两个文件, 如果 example文件夹不存在, 则会创建
#   baseName: example/server
# # cfg 是证书的主要配置
#   cfg:
#     # algorithm 的取值关系是
#     # 1----RSA,
#     # 2----DSA,
#     # 3----ECDSA,
#     # 4----Ed25519,
#     # 目前只支持 1,3 即 RSA 和 ECSDA
#     algorithm: 1
#     # 证书有效期, 年为单位
#     duration: 99
#     # 证书用途，openssl v3 扩展
#     # 1--ExtKeyUsageServerAuth
#     # 2--ExtKeyUsageClientAuth
#     # 3--ExtKeyUsageCodeSigning
#     # 4--ExtKeyUsageEmailProtection
#     # 5--ExtKeyUsageIPSECEndSystem
#     # 6--ExtKeyUsageIPSECTunnel
#     # 7--ExtKeyUsageIPSECUser
#     # 8--ExtKeyUsageTimeStamping
#     # 9--ExtKeyUsageOCSPSigning
#     # 10-ExtKeyUsageMicrosoftServerGatedCrypto
#     # 11-ExtKeyUsageNetscapeServerGatedCrypto
#     # 12-ExtKeyUsageMicrosoftCommercialCodeSigning
#     # 13-ExtKeyUsageMicrosoftKernelCodeSigning
#     usages:
#       - 1
#       - 2
#     # x509 SAN 扩展, dns 和 ip 均可写多个
#     altNames:
#       dns:
#         - localhost
#       ips:
#         - 127.0.0.1
#     # 证书 CN
#     commonName: example-server
#     # 证书 OU, 可写多个
#     organization:
#       - "system:master"

- name: etcd-ca
  baseName: etcd/ca
  caName: ""
  cfg:
    algorithm: 1
    commonName: etcd-ca
    # In years, ignoring leap years， 1 year = 365 day
    duration: 99
- name: etcd-server
  baseName: etcd/server
  caName: "etcd-ca"
  cfg:
    algorithm: 1
    duration: 99
    usages:
      - 1
      - 2
    altNames:
      dns:
        - localhost
      ips:
        - 127.0.0.1
    commonName: etcd-server
    organization:
      - "system:master"
- name: etcd-peer
  baseName: etcd/peer
  caName: "etcd-ca"
  cfg:
    algorithm: 1
    duration: 99
    usages:
      - 1
      - 2
    altNames:
      dns:
        - localhost
      ips:
        - 127.0.0.1
    commonName: etcd-peer
    organization:
      - "system:master"
- name: etcd-healthcheck-client
  baseName: etcd/healthcheck-client
  caName: "etcd-ca"
  cfg:
    algorithm: 1
    duration: 99
    usages:
      - 2
    commonName: etcd-healthcheck-client
    organization:
      - "system:master"
`
)
