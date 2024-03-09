import { Fetcher, KVNamespace } from "@cloudflare/workers-types";

export interface Env {
  AI: Fetcher;
  SECTIONS: KVNamespace;
  SECRET_KEY: string;
}
