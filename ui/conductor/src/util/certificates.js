import * as asn1js from 'asn1js';
import * as pkijs from 'pkijs';

/**
 * Convert a binary string to an ArrayBuffer.
 *
 * @param {string} binaryString
 * @return {ArrayBuffer}
 */
function binaryStringToArrayBuffer(binaryString) {
  const len = binaryString.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }

  return bytes.buffer;
}

/**
 * Convert an ArrayBuffer to a hex string.
 *
 * @param {ArrayBuffer} arrayBuffer
 * @return {string}
 */
function arrayBufferToHex(arrayBuffer) {
  const hexString = Array.from(new Uint8Array(arrayBuffer))
      .map(byte => byte.toString(16).padStart(2, '0'))
      .join('').toUpperCase();

  return hexString.replace(/^0+/, '');
}

// Format a hex string with colons between each pair of digits.
/**
 *
 * @param {string} hexString
 * @return {string}
 */
function formatFingerprint(hexString) {
  return hexString.match(/.{1,2}/g).join(':');
}

/**
 * Get the algorithm name from an OID.
 *
 * @param {string} algorithmOid
 * @return {string}
 */
function getAlgorithmName(algorithmOid) {
  const oidMap = {
    '1.2.840.113549.1.1.1': 'RSA Encryption',
    '1.2.840.113549.1.1.11': 'SHA-256 with RSA Encryption',
    '1.2.840.113549.1.1.5': 'SHA-1 with RSA Encryption',
    '1.2.840.10040.4.1': 'DSA',
    '1.2.840.10045.2.1': 'EC Public Key',
    '1.2.840.10045.4.3.2': 'SHA-256 with ECDSA',
    '1.2.840.10045.4.3.3': 'SHA-384 with ECDSA',
    '1.2.840.10045.4.3.4': 'SHA-512 with ECDSA',
    '1.2.840.113549.1.1.4': 'MD5 with RSA Encryption',
    '1.2.840.113549.1.1.2': 'MD2 with RSA Encryption',
    '1.3.14.3.2.26': 'SHA-1',
    '1.3.14.3.2.29': 'SHA-1 with RSA Signature',
    '1.2.840.113549.1.1.10': 'RSA-PSS',
    '1.2.840.113549.1.1.12': 'SHA-256 with RSA and MGF1',
    '1.2.840.113549.1.1.13': 'SHA-384 with RSA and MGF1',
    '1.2.840.113549.1.1.14': 'SHA-512 with RSA and MGF1'
    // ... Add other OID mappings as needed ...
  };

  return oidMap[algorithmOid] || 'Unknown Algorithm';
}

/**
 * Calculate SHA1 hash of the given ArrayBuffer.
 *
 * @param {ArrayBuffer} buffer
 * @param {string} algorithm
 * @return {Promise<string>} Hexadecimal hash string.
 */
async function shaHex(buffer, algorithm) {
  const hashBuffer = await window.crypto.subtle.digest(algorithm, buffer);
  return arrayBufferToHex(hashBuffer);
}


// ---------------------------- //
/**
 * Parse a PEM certificate and return an ArrayBuffer representation.
 *
 * @param {string} pem - The PEM string of the certificate.
 * @return {Promise<{
 * validityPeriod: string,
 * extensions: {keyUsage: string, basicConstraints: string, subjectKeyIdentifier: string},
 * keyLength: number,
 * serial: string,
 * subject: {commonName: string, organization: string},
 * sha1Fingerprint: string,
 * sha256Fingerprint: string,
 * primaryDomain: string,
 * subjectAltDomains: string,
 * version: number,
 * signatureAlgorithm: string,
 * issuer: {commonName: string}
 * }>} - The ArrayBuffer representation of the certificate.
 */
export async function parseCertificate(pem) {
  // Define PEM header and footer
  const pemHeader = '-----BEGIN CERTIFICATE-----';
  const pemFooter = '-----END CERTIFICATE-----';

  // Remove the header, footer, and newlines
  const pemContents = pem.replace(pemHeader, '').replace(pemFooter, '').replace(/\r\n|\n|\r/gm, '');

  // Ensure that pemContents is properly Base64-encoded
  if (!pemContents.match(/^[A-Za-z0-9+/=]+$/)) {
    throw new Error('Invalid PEM encoded string.');
  }

  // Decode Base64 to binary string
  const binaryString = window.atob(pemContents);
  const arrayBuffer = binaryStringToArrayBuffer(binaryString);

  // Use ASN.1 parser
  const asn1 = asn1js.fromBER(arrayBuffer);
  const cert = new pkijs.Certificate({schema: asn1.result});

  const certificateInformation = {
    primaryDomain: '',
    subjectAltDomains: '',
    validityPeriod: {
      from: '',
      to: ''
    },
    signatureAlgorithm: '',
    keyLength: 0,
    serial: '',
    sha1Fingerprint: '',
    sha256Fingerprint: '',
    version: 0,
    subject: {
      organization: '',
      commonName: ''
    },
    issuer: {
      commonName: ''
    },
    extensions: {
      keyUsage: '',
      basicConstraints: '',
      subjectKeyIdentifier: ''
    }
  };

  // Subject & Primary Domain
  if (cert.subject && cert.subject.typesAndValues) {
    const primaryDomainValue = cert.subject.typesAndValues.find(t => t.type === '2.5.4.3'); // Common Name (CN)
    certificateInformation.primaryDomain = primaryDomainValue ? primaryDomainValue.value.valueBlock.value : '';

    const orgValue = cert.subject.typesAndValues.find(t => t.type === '2.5.4.10'); // Organization (O)
    certificateInformation.subject.organization = orgValue ? orgValue.value.valueBlock.value : '';

    const cnValue = cert.subject.typesAndValues.find(t => t.type === '2.5.4.3'); // Common Name (CN)
    certificateInformation.subject.commonName = cnValue ? cnValue.value.valueBlock.value : '';
  }

  // Issuer
  if (cert.issuer && cert.issuer.typesAndValues) {
    const issuerCNValue = cert.issuer.typesAndValues.find(t => t.type === '2.5.4.3'); // Common Name (CN)
    certificateInformation.issuer.commonName = issuerCNValue ? issuerCNValue.value.valueBlock.value : '';
  }

  // Validity period
  const dateTimeOptions = {
    weekday: 'short',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  };

  if (cert.notBefore && cert.notBefore.value) {
    const notBeforeDate = new Date(cert.notBefore.value);
    certificateInformation.validityPeriod.from = notBeforeDate.toLocaleString(undefined, dateTimeOptions);
  }

  if (cert.notAfter && cert.notAfter.value) {
    const notAfterDate = new Date(cert.notAfter.value);
    certificateInformation.validityPeriod.to = notAfterDate.toLocaleString(undefined, dateTimeOptions);
  }

  // Signature algorithm
  if (cert.signatureAlgorithm && cert.signatureAlgorithm.algorithmId) {
    certificateInformation.signatureAlgorithm = getAlgorithmName(cert.signatureAlgorithm.algorithmId);
  }

  // Key length
  // if (cert.subjectPublicKeyInfo && cert.subjectPublicKeyInfo.parsedKey) {
  //   const publicKey = cert.subjectPublicKeyInfo.parsedKey;
  //   if (publicKey.modulus) { // For RSA keys
  //     const modulusHex = publicKey.modulus.valueBlock.valueHex;
  //     certificateInformation.keyLength = modulusHex.byteLength * 8;
  //   }
  // }

  // Serial number
  // if (cert.serialNumber) {
  //   certificateInformation.serial = arrayBufferToHex(cert.serialNumber.valueBlock.valueHex);
  // }

  // Version
  if (cert.version) {
    certificateInformation.version = cert.version;
  }

  // Extensions
  if (cert.extensions) {
    // Key Usage
    const keyUsageExtension = cert.extensions.find(ext => ext.extnID === '2.5.29.15');
    if (keyUsageExtension) {
      const keyUsageValueHex = keyUsageExtension.parsedValue.valueBlock.valueHex;
      const keyUsageBuffer = new Uint8Array(keyUsageValueHex);
      const keyUsages = [];

      // Bit positions for key usages (as per X.509 specification)
      if ((keyUsageBuffer[0] & 0x80) !== 0) keyUsages.push('Digital Signature');
      if ((keyUsageBuffer[0] & 0x40) !== 0) keyUsages.push('Non Repudiation');
      if ((keyUsageBuffer[0] & 0x20) !== 0) keyUsages.push('Key Encipherment');
      if ((keyUsageBuffer[0] & 0x10) !== 0) keyUsages.push('Data Encipherment');
      if ((keyUsageBuffer[0] & 0x08) !== 0) keyUsages.push('Key Agreement');
      if ((keyUsageBuffer[0] & 0x04) !== 0) keyUsages.push('Certificate Sign');
      if ((keyUsageBuffer[0] & 0x02) !== 0) keyUsages.push('CRL Sign');
      if ((keyUsageBuffer[0] & 0x01) !== 0) keyUsages.push('Encipher Only');
      if (keyUsageBuffer.length > 1 && (keyUsageBuffer[1] & 0x80) !== 0) keyUsages.push('Decipher Only');

      certificateInformation.extensions.keyUsage = keyUsages.join(', ');
    }

    // Basic Constraints
    const basicConstraintsExtension = cert.extensions.find(ext => ext.extnID === '2.5.29.19');
    if (basicConstraintsExtension) {
      certificateInformation.extensions.basicConstraints = `CA: ${basicConstraintsExtension.parsedValue.cA}`;
    }

    // Subject Key Identifier
    const subjectKeyIdentifierExtension = cert.extensions.find(ext => ext.extnID === '2.5.29.14');
    if (subjectKeyIdentifierExtension) {
      const ski = arrayBufferToHex(subjectKeyIdentifierExtension.parsedValue.valueBlock.valueHex);
      certificateInformation.extensions.subjectKeyIdentifier = formatFingerprint(ski);
    }
  }

  // SHA1 Fingerprint
  certificateInformation.sha1Fingerprint = formatFingerprint(await shaHex(arrayBuffer, 'SHA-1'));

  // SHA256 Fingerprint
  certificateInformation.sha256Fingerprint = formatFingerprint(await shaHex(arrayBuffer, 'SHA-256'));

  return certificateInformation;
}
