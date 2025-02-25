import axios, { AxiosResponse, AxiosRequestConfig } from "axios";
import { useStore } from "@/state/store";

export enum APIMethod {
  GET = "get",
  POST = "post",
  UPDATE = "put",
  DELETE = "delete",
  POST_FILE = "postFile",
}

export interface APIError {
  error: true;
  status: number;
  message: string;
}

export const request = async (
  body: any,
  endpoint: string,
  method: APIMethod
): Promise<AxiosResponse | APIError> => {
  const headers: Record<string, string> = {
    "Content-Type":
      method === APIMethod.POST_FILE
        ? "multipart/form-data"
        : "application/json",
    Network: useStore.getState().network,
  };

  const config: AxiosRequestConfig = {
    headers,
    method: method === APIMethod.UPDATE ? "put" : method,
    url: endpoint,
  };

  if (method === APIMethod.GET) {
    config.params = body;
  } else if (method === APIMethod.POST_FILE) {
    const form = new FormData();
    form.append("file", body.Filename);
    config.data = form;
  } else {
    config.data = body;
  }

  try {
    const response = await axios(config);
    return response;
  } catch (error: any) {
    console.error("E_REQUEST =>", error);
    return {
      error: true,
      status: error.response?.status || 500,
      message:
        error.response?.data?.message ||
        error.message ||
        "Internal Server Error",
    };
  }
};
