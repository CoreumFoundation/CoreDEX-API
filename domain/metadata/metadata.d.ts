import _m0 from "protobufjs/minimal";
export declare const protobufPackage = "metadata";
export declare enum Network {
    /** NETWORK_DO_NOT_USE - Avoid a default value since a default and a bug are not distinguishable */
    NETWORK_DO_NOT_USE = 0,
    MAINNET = 1,
    TESTNET = 2,
    DEVNET = 3,
    UNRECOGNIZED = -1
}
export declare function networkFromJSON(object: any): Network;
export declare function networkToJSON(object: Network): string;
export interface MetaData {
    /** mainnet, testnet, devnet, can also be some virtually defined network (extra devnet for testing, extra mainnnet node for scanning historical blocks, etc) */
    Network: Network;
    UpdatedAt: Date | undefined;
    CreatedAt: Date | undefined;
}
export declare const MetaData: {
    encode(message: MetaData, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MetaData;
    fromJSON(object: any): MetaData;
    toJSON(message: MetaData): unknown;
    create<I extends {
        Network?: Network | undefined;
        UpdatedAt?: Date | undefined;
        CreatedAt?: Date | undefined;
    } & {
        Network?: Network | undefined;
        UpdatedAt?: Date | undefined;
        CreatedAt?: Date | undefined;
    } & { [K in Exclude<keyof I, keyof MetaData>]: never; }>(base?: I | undefined): MetaData;
    fromPartial<I_1 extends {
        Network?: Network | undefined;
        UpdatedAt?: Date | undefined;
        CreatedAt?: Date | undefined;
    } & {
        Network?: Network | undefined;
        UpdatedAt?: Date | undefined;
        CreatedAt?: Date | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MetaData>]: never; }>(object: I_1): MetaData;
};
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P : P & {
    [K in keyof P]: Exact<P[K], I[K]>;
} & {
    [K in Exclude<keyof I, KeysOfUnion<P>>]: never;
};
export {};
