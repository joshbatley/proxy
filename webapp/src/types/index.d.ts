export type Collections = {
  id: string;
  name: string
  endpoints: Endpoints[]
};

export type Endpoints = {
  id: string;
  status: number;
  url: string;
  method: string;
};
