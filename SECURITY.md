# Security Policy

## Supported Versions

The following versions of ads-go are currently supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.2.x   | :white_check_mark: |
| 0.1.x   | :x:                |

We recommend always using the latest release for the best security and stability.

## Reporting a Vulnerability

We take the security of ads-go seriously. If you discover a security vulnerability, please follow these steps:

### 1. **Do Not** Create a Public Issue

Please do not open a public GitHub issue for security vulnerabilities, as this could put users at risk.

### 2. Report Privately

Report security vulnerabilities through one of these methods:

#### GitHub Security Advisories (Preferred)

1. Go to the [Security tab](https://github.com/JarmoCluyse/ads-go/security) of the repository
2. Click "Report a vulnerability"
3. Fill out the form with details about the vulnerability

#### Email

Alternatively, email security concerns to: **jarmo_cluyse@hotmail.com**

### 3. What to Include

When reporting a vulnerability, please include:

- **Description**: Clear description of the vulnerability
- **Impact**: What could an attacker do with this vulnerability?
- **Reproduction**: Step-by-step instructions to reproduce the issue
- **Affected versions**: Which versions of ads-go are affected?
- **Environment**: Go version, OS, TwinCAT version (if relevant)
- **Proof of concept**: Code demonstrating the vulnerability (if possible)
- **Suggested fix**: If you have ideas for how to fix it

### 4. What to Expect

After you report a vulnerability:

- **Initial response**: Within 48 hours, we'll acknowledge receipt
- **Assessment**: Within 1 week, we'll assess the vulnerability and provide initial feedback
- **Updates**: We'll keep you informed of progress every 1-2 weeks
- **Fix timeline**: Critical vulnerabilities will be addressed within 30 days
- **Disclosure**: We'll coordinate with you on responsible disclosure timing

## Security Best Practices

When using ads-go in your applications:

### Network Security

- **Firewall**: Restrict ADS port access (TCP 48898) to trusted networks
- **TLS/VPN**: Use TLS or VPN when connecting over untrusted networks
- **Authentication**: Configure TwinCAT router authentication where possible
- **Network isolation**: Run PLCs and ADS clients on isolated networks

### Credentials

- **No hardcoding**: Never hardcode AMS NetIDs, IP addresses, or credentials in source code
- **Environment variables**: Use environment variables or secure configuration files
- **Secrets management**: Use proper secrets management for production deployments
- **Least privilege**: Grant only necessary ADS permissions

### Input Validation

- **Sanitize input**: Validate all user input before sending to PLC
- **Type checking**: Use strong typing and validate data types match PLC expectations
- **Bounds checking**: Check array indices and string lengths before writing
- **Error handling**: Always handle errors from ads-go operations

### Dependencies

- **Keep updated**: Regularly update ads-go to the latest version
- **Dependency scanning**: Use tools like `go list -m all` and vulnerability scanners
- **Vendor dependencies**: Consider vendoring dependencies for production

### Logging

- **Sanitize logs**: Don't log sensitive data (credentials, production values)
- **Rate limiting**: Implement rate limiting on ADS operations
- **Monitoring**: Monitor for unusual patterns or errors

## Known Security Considerations

### ADS Protocol Limitations

The ADS protocol has inherent security limitations:

- **No encryption**: ADS does not encrypt data in transit (use VPN/TLS at network layer)
- **Limited authentication**: Authentication is router-based, not per-connection
- **Industrial protocol**: Designed for trusted industrial networks, not internet-facing applications

### Recommendations

1. **Never expose ADS directly to the internet**
2. **Use network segmentation** to isolate PLCs
3. **Implement application-layer security** (authentication, authorization, rate limiting)
4. **Monitor and audit** all ADS operations in production

## Security Updates

Security updates will be:

- Released as patch versions (e.g., 0.2.1)
- Documented in [CHANGELOG.md](CHANGELOG.md)
- Announced in GitHub releases
- Tagged with "security" label

Critical vulnerabilities may result in out-of-band releases.

## Disclosure Policy

When a vulnerability is fixed:

1. **Private fix**: We'll develop and test the fix privately
2. **Security advisory**: We'll publish a GitHub Security Advisory
3. **Release**: We'll release a new version with the fix
4. **Announcement**: We'll announce the security update in the release notes
5. **Credit**: We'll credit the reporter (unless they prefer anonymity)

We follow a **90-day disclosure policy**: vulnerabilities will be publicly disclosed 90 days after the fix is released, or when the vulnerability becomes publicly known, whichever comes first.

## Contact

For security concerns: **jarmo_cluyse@hotmail.com**

For general questions: Open a [GitHub Discussion](https://github.com/JarmoCluyse/ads-go/discussions)

---

Thank you for helping keep ads-go and its users secure!
