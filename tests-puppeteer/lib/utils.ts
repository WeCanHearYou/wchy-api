import "reflect-metadata";
import { Browser } from ".";

export const delay = (ms: number) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

export function findBy(selector: string) {
  return (target: any, propertyKey: string) => {
    const type = Reflect.getMetadata("design:type", target, propertyKey);
    Object.defineProperty(target, propertyKey, {
      configurable: true,
      enumerable: true,
      get() {
        const browser = (this as any).browser;
        const handle = (browser as Browser).findElement(selector);
        return new type(handle, selector, browser);
      }
    });
  };
}

// export function findMultipleBy(selector: string, t?: new (element: WebElementPromise, selector: string) => any) {
//   return (target: any, propertyKey: string) => {
//     const type = Reflect.getMetadata("design:type", target, propertyKey);
//     Object.defineProperty(target, propertyKey, {
//       configurable: true,
//       enumerable: true,
//       get() {
//         const browser = (this as any).browser;
//         const promise = (browser as Browser).findElements(selector);
//         return new type(promise, selector, browser);
//       }
//     });
//   };
// }
