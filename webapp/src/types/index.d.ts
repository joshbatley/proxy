export type Collections = {
  id: string;
  name: string
  endpoints: Endpoint[]
};

export type Endpoint = {
  id: string;
  status: number;
  url: string;
  method: string;
};

export type Response = {
  id:  string;
  status:  number;
  url:  string;
  method: string;
  headers: string;
  body: string;
  datetime: number;
}
