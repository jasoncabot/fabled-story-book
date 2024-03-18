interface Source {
  id: string;
  name: string;
  url: string;
  entrypoint: string;
  auth: () => Promise<string | undefined>;
}

export const defaultEntrypoint = "entrpoint.jabl";

export const sources: Source[] = [
  {
    id: "1",
    name: "Example 1",
    url: "https://raw.githubusercontent.com/jasoncabot/fabled-story-book/main/assets/example01/",
    entrypoint: "entrypoint.jabl",
    auth: () => {
      return Promise.resolve(undefined);
    },
  },
  {
    id: "2",
    name: "1. The War-Torn Kingdom",
    url: `/story/1-war-torn-kingdom/`,
    entrypoint: "0-choose-character.jabl",
    auth: () => {
      return Promise.resolve(undefined);
    },
  },
];

export const availableSources = sources.reduce((acc: Record<string, string>, source) => {
  acc[source.id] = source.url;
  return acc;
}, {});

export const authFunctions = sources.reduce((acc: Record<string, () => Promise<string | undefined>>, source) => {
  acc[source.id] = source.auth;
  return acc;
}, {});
