declare module "*.wasm" {
  const content: string;
  export default content;
}

declare class Go {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
}

declare global {
  interface Window {
    execSection: (section: string, callback: (code: string, error: string) => void) => void;
    evalCode: (code: string, callback: (result: string, error: string) => void) => void;
  }
}
