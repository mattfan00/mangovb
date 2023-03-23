const base = async <T>(url: RequestInfo | URL, options: RequestInit)=> {
    const res = await fetch(url, options);
    return await res.json() as T;
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
