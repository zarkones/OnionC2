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

export type Operator = {
    username: string
    publicKeyHex: string
    createdAt: number
}

export const API = ref(new class {
    public username: string
    public c2HostURL: string
    
    private privateKey: CryptoKey | null
    private store
    
    constructor() {
        this.username = ''
        this.c2HostURL = ''
        this.privateKey = null
        this.store = {
            operators: {
                data: [] as Operator[],
                page: 0,
            },
            agents: {
                data: [] as Agent[],
                page: 0,
            },
        }

        this.initializePeriodicDataFetching()
    }

    private initializePeriodicDataFetching = () => {
        setInterval(async () => {
            await Promise.all([
                async () => this.store.agents.data = await this.fetchAgents(this.store.agents.page),
                async () => this.store.operators.data = await this.fetchOperators(this.store.operators.page),
            ])
        }, 14000)
    }

    public getAgents = () => this.store.agents

    public getOperators = () => this.store.operators

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
            .setExpirationTime('3min')
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

    private fetchAgents = async (page: number) => {
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

    private fetchOperators = async (page: number) => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        } 

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/operators?page=${page}`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return []
        }

        return await response.json() as Operator[]
    }
}())