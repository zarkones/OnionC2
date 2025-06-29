import * as jose from 'jose'
import { ref } from 'vue'

export type JWTPayload = {
    u: string // Operator's username.
}

export type Agent = {
    id: string
    hostname: string
    os: string
    ip: string
    ram: string
    osVersion: string
    cpuName: string
    arch: string
    lastSeen: number
}

export const API = ref(new class {
    public username: string
    public c2HostURL: string
    
    private privateKey: CryptoKey | null
    
    constructor() {
        this.username = ''
        this.c2HostURL = ''
        this.privateKey = null
    }

    public setPrivateKey = async (pemEncodedPrivateKey: string) => {
        // Converting from older PKCS#1 to PKCSSss8 header.
        pemEncodedPrivateKey = pemEncodedPrivateKey.replaceAll('-----BEGIN RSA PRIVATE KEY-----', '-----BEGIN PRIVATE KEY-----')
        pemEncodedPrivateKey = pemEncodedPrivateKey.replaceAll('-----END RSA PRIVATE KEY-----', '-----END PRIVATE KEY-----')

        this.privateKey = await jose.importPKCS8(pemEncodedPrivateKey, 'RS512')
    } 

    public sign = async (data: any) => {
        if (this.privateKey === null) {
            throw Error('private signing key is null')
        }
        const jwt = await new jose.SignJWT(data)
            .setProtectedHeader({ alg: 'RS512' })
            .setExpirationTime('1d')
            .setNotBefore(0)
            .sign(this.privateKey)
        return jwt
    }

    public verify = async (verificationKey: string, token: string) => {
        const publicKey = await jose.importSPKI(verificationKey, 'RS512')
        const { payload } = await jose.jwtVerify(token, publicKey, {
            algorithms: ['RS512']
        })
        return payload
    }

    public getAgents = async (page: number) => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        } 

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/agents?page=${page}`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return []
        }

        return await response.json() as Agent[]
    }
}())