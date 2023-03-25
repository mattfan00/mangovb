export interface Res<T> {
    data?: T;
    error?: Error;
}

interface ErrResponse {
    message: string;
}

const base = async <T>(url: RequestInfo | URL, options: RequestInit): Promise<Res<T>> => {
    const res = await fetch(url, options);
    if (!res.ok) {
        const data = await res.json() as ErrResponse;

        return {
            error: new Error(data.message),
        };
    }

    return {
        data: await res.json() as T,
    };
};

export const get = async <T>(url: string) => {
    return await base<T>(url, {});
};

export const withSearchParams = (url: string, searchParams: string | URLSearchParams) => {
    let u = new URL(url);
    u.search = "";
    const s = new URLSearchParams(searchParams);
    if (Array.from(s.entries()).length > 0) {
        u = new URL(`${u.origin}${u.pathname}?${new URLSearchParams(searchParams)}`); 
    }

    return u.href;
};

export default {
    get,
    withSearchParams,
};
