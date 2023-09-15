import qs from 'qs';

export type Product = {
    id: string;
    title: string;
    preview_img: string;
    description: string;
    price: string;
    url: string;
};

export type FetchAllParameters = {
    query: {
        query: string | undefined;
        sources: string[];
        categories: string[];
        price: {
            from: number | undefined;
            to: number | undefined;
        };
    };
};

export default function () {
    const { public: publicConfig } = useRuntimeConfig();

    const fetchAll = (params: FetchAllParameters): Promise<Product[]> => {
        return $fetch<Product[]>('/search', {
            baseURL: publicConfig.api.baseUrl || 'http://localhost',
            params: {
                product: params.query.query,
                sources: params.query.sources,
                categories: params.query.categories,
                price: {
                    from: params.query.price.from,
                    to: params.query.price.to,
                },
            },
            onRequest: (ctx) => {
                if (ctx.options.params || ctx.options.query) {
                    ctx.request = ctx.request + qs.stringify({
                        ...ctx.options.query,
                        ...ctx.options.params,
                    }, {
                        arrayFormat: 'brackets',
                        addQueryPrefix: true,
                        skipNulls: true,
                    });

                    ctx.options.params = undefined;
                    ctx.options.query = undefined;
                }
            },
            headers: {
                'accept': 'application/json',
            },
            // TODO: delete this when backend is ready
            parseResponse(string) {
                return JSON.parse(string);
            },
        });
    };
    
    const fetchOne = () => {
        // TODO: implement when the backend is ready
    };

    return {
        fetchAll,
        fetchOne,
    };
}