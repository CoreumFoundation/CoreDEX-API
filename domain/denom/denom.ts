// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v1.181.2
//   protoc               v3.20.0
// source: domain/denom/denom.proto

/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "denom";

export interface Denom {
  Currency: string;
  Issuer: string;
  Precision?: number | undefined;
  IsIBC: boolean;
  /** Could be handy for IBC */
  Denom: string;
  /** Additional fields, make it possible to use the denom as the currency storage for display purposes */
  Name?: string | undefined;
  Description?: string | undefined;
  Icon?: string | undefined;
}

function createBaseDenom(): Denom {
  return {
    Currency: "",
    Issuer: "",
    Precision: undefined,
    IsIBC: false,
    Denom: "",
    Name: undefined,
    Description: undefined,
    Icon: undefined,
  };
}

export const Denom = {
  encode(message: Denom, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Currency !== "") {
      writer.uint32(10).string(message.Currency);
    }
    if (message.Issuer !== "") {
      writer.uint32(18).string(message.Issuer);
    }
    if (message.Precision !== undefined) {
      writer.uint32(24).int32(message.Precision);
    }
    if (message.IsIBC !== false) {
      writer.uint32(32).bool(message.IsIBC);
    }
    if (message.Denom !== "") {
      writer.uint32(42).string(message.Denom);
    }
    if (message.Name !== undefined) {
      writer.uint32(50).string(message.Name);
    }
    if (message.Description !== undefined) {
      writer.uint32(58).string(message.Description);
    }
    if (message.Icon !== undefined) {
      writer.uint32(66).string(message.Icon);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Denom {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDenom();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Currency = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Issuer = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.Precision = reader.int32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.IsIBC = reader.bool();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.Denom = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.Name = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.Description = reader.string();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.Icon = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Denom {
    return {
      Currency: isSet(object.Currency) ? globalThis.String(object.Currency) : "",
      Issuer: isSet(object.Issuer) ? globalThis.String(object.Issuer) : "",
      Precision: isSet(object.Precision) ? globalThis.Number(object.Precision) : undefined,
      IsIBC: isSet(object.IsIBC) ? globalThis.Boolean(object.IsIBC) : false,
      Denom: isSet(object.Denom) ? globalThis.String(object.Denom) : "",
      Name: isSet(object.Name) ? globalThis.String(object.Name) : undefined,
      Description: isSet(object.Description) ? globalThis.String(object.Description) : undefined,
      Icon: isSet(object.Icon) ? globalThis.String(object.Icon) : undefined,
    };
  },

  toJSON(message: Denom): unknown {
    const obj: any = {};
    if (message.Currency !== "") {
      obj.Currency = message.Currency;
    }
    if (message.Issuer !== "") {
      obj.Issuer = message.Issuer;
    }
    if (message.Precision !== undefined) {
      obj.Precision = Math.round(message.Precision);
    }
    if (message.IsIBC !== false) {
      obj.IsIBC = message.IsIBC;
    }
    if (message.Denom !== "") {
      obj.Denom = message.Denom;
    }
    if (message.Name !== undefined) {
      obj.Name = message.Name;
    }
    if (message.Description !== undefined) {
      obj.Description = message.Description;
    }
    if (message.Icon !== undefined) {
      obj.Icon = message.Icon;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Denom>, I>>(base?: I): Denom {
    return Denom.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Denom>, I>>(object: I): Denom {
    const message = createBaseDenom();
    message.Currency = object.Currency ?? "";
    message.Issuer = object.Issuer ?? "";
    message.Precision = object.Precision ?? undefined;
    message.IsIBC = object.IsIBC ?? false;
    message.Denom = object.Denom ?? "";
    message.Name = object.Name ?? undefined;
    message.Description = object.Description ?? undefined;
    message.Icon = object.Icon ?? undefined;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
