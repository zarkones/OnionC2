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

export type Message = {
    id: string
    agentId: string
    request: string
    response: string
    createdAt: number
}

export const API = ref(new class {
    public username: string
    public c2HostURL: string
    
    private privateKey: CryptoKey | null
    public store
    
    constructor() {
        this.username = 'zarkones3'
        this.c2HostURL = 'http://localhost:8080'
        this.privateKey = null
        this.store = ref({
            operators: {
                data:[] as Operator[],
                page: 0,
            },
            agents: {
                data:[] as Agent[],
                page: 0,
            },
            messages: {
                newMessageCallback: () => {},
                data:[] as Message[],
                agentId: '',
                lastSentAt: 0,
                page: undefined as number | undefined,
                since: undefined as string | undefined,
            },
        })

        this.initializePeriodicDataFetching()
    }

    private initializePeriodicDataFetching = async () => {
        const generalUpdate = async () => {
            if (this.privateKey === null || this.c2HostURL.length === 0 || this.username.length === 0) {
                return
            }

            this.store.value.agents.data = await this.fetchAgents(this.store.value.agents.page)
            this.store.value.operators.data = await this.fetchOperators(this.store.value.operators.page)
        }

        await generalUpdate()

        setInterval(async () => {
            await generalUpdate()
        }, 3000)
    }

    public clearMessages = () => {
        this.store.value.messages = {
            newMessageCallback: () => {},
            lastSentAt: 0,
            agentId: '',
            data: [],
            page: undefined,
            since: undefined,
        }
    }

    public setPrivateKey = async (pemEncodedPrivateKey: string) => {
        // Converting from older PKCS#1 to PKCSS#8 header.
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

    public sendMessage = async (agentId: string, content: string) => {
        const body: Partial<Message> = {
            agentId,
            request: content,
        }

        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/messages`, {
            method: 'PUT',
            headers: {
                Authorization: token,
            },
            body: JSON.stringify(body),
        })

        if (response.status !== 201) {
            throw Error(`unexpected status code: ${response.statusText}`)
        }

        this.store.value.messages.lastSentAt = Date.now()
    }

    public fetchMessages = async (agentId: string, options: { page: number | undefined, since: string | undefined }) => {
        if (agentId.length === 0) {
            return { since: "", messages: [] }
        }
        
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)

        const queryParams = [
            typeof options.page !== undefined ? `page=${options.page}` : '',
            typeof options.since !== undefined ? `since=${options.since}` : '',
        ]

        const response = await fetch(`${this.c2HostURL}/v1/messages/${agentId}?${queryParams.join('&')}`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return { since: "", messages: [] }
        }

        const r = await response.json() as { since: string, messages: Message[] } 
        
        return { ...r, messages: r.messages.reverse() }
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