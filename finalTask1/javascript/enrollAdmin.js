/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'connection-Manufacturer.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

async function main() {
    try {

        console.log("File Path --> ",__dirname);
            // Create a new CA client fors interacting with the CA.
        const caInfo = ccp.certificateAuthorities['ca.Manufacturer.example.com'];
        console.log("caInfo --> ",caInfo)
        const caTLSCACertsPath = path.resolve(__dirname, '..','..','finalTask', caInfo.tlsCACerts.path);
        console.log("caTLSCACertsPath Path --> ",caTLSCACertsPath);

        const caTLSCACerts = fs.readFileSync(caTLSCACertsPath);
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);
        console.log("CA --> ",ca);


        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists('admin');
        if (adminExists) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
        console.log("enrollment --> ",enrollment);

        const identity = X509WalletMixin.createIdentity('ManufacturerMSP', enrollment.certificate, enrollment.key.toBytes());
        console.log("identity --> ",identity);

        await wallet.import('admin', identity);

        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');

    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
        process.exit(1);
    }
}

main();
