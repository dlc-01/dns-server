# DNS Server (Educational)

An educational implementation of a DNS-compatible server built from scratch.  
The project focuses on understanding the DNS protocol by implementing packet parsing, name compression, and basic request handling over UDP.

This is not a production-ready DNS server.

---

## Implemented Features

- UDP server for handling DNS requests
- DNS wire protocol parsing and response building
- Header, Question, and Answer sections support
- DNS name encoding and compression (pointers)
- Forwarding queries to an upstream DNS resolver
- Correct transaction ID and flag handling

---

## Supported Scenarios

### DNS Queries

- Single and multiple questions per request
- Correct parsing of QNAME, QTYPE, QCLASS
- Building valid DNS responses

### Answers

- Single and multiple answers
- Proper handling of NAME, TYPE, CLASS, TTL, RDLENGTH, RDATA

### Name Compression

- Parsing compressed domain names
- Handling pointer-based jumps and offsets

### Forwarding

- Forward queries to an upstream resolver over UDP
- Read and merge upstream responses

### Errors

- Unsupported opcode
- Resolver failures
- Best-effort handling of malformed packets

---

## Limitations

- UDP only (no TCP fallback)
- No caching
- Limited record type support
- No EDNS or DNSSEC

---

## Purpose

This project is intended for learning DNS internals, binary protocols, and clean server architecture.
