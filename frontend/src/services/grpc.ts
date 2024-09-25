import {AuthService as originAuthService} from "@/lib/proto/auth/v1/auth.pb.ts";
import {UserService as originUserService} from "@/lib/proto/user/v1/user.pb.ts";
import {CategoryService as originCategoryService} from "@/lib/proto/category/v1/category.pb.ts";
import {useUserStore} from "@/store";

export const AuthService = createProxy(originAuthService);
export const UserService = createProxy(originUserService);
export const CategoryService = createProxy(originCategoryService);

function createProxy<T>(service: T): T {
	const handler:ProxyHandler<any> = {
		get: function (target: any, prop: any) {
			if (target[prop] && (typeof target[prop]) === "function") {
				return async (...args: any[]) => {
					if (args.length == 0) {
						args.push({});
					}

					const opts = args[args.length - 1];
					// check if the last arg is a function (streaming request)
					// or if has only a single argument
					if ((typeof opts) === "function" || args.length == 1) {
						const prefix = {
							pathPrefix: import.meta.env.VITE_GAPI_URL,
							headers: {} as Record<string, string>,
						};
						const token = useUserStore().token;
						if (token){
							prefix.headers["Authorization"] = `Bearer ${token}`
						}else {
							delete prefix.headers["Authorization"]
						}

						args.push(prefix);
					}

					return await target[prop](...args).catch(async (err: RequestError) => {
						err = new RequestError(err)
						if (err.code == 16) {
							try {
								await useUserStore().refreshToken().then(async () => {
									const newToken = useUserStore().token;
									if (newToken) {
										args[args.length - 1].headers["Authorization"] = `Bearer ${newToken}`;
									}
									return await target[prop](...args);
								})
							} catch (err) {
								useUserStore().$reset();
								if (window.location.pathname !== '/login') {
									const query = new URLSearchParams();
									query.set('from', window.location.pathname + window.location.search);
									window.location.href = window.location.origin + "/login?" + query.toString()
								}
							}
						}
						throw err
					});
				};
			}

			return target[prop];
		}
	}
	return new Proxy(service, handler)
}

export class RequestError extends Error {
	public code = 0;
	public message = "";

	constructor(response: any) {
		super();
		if (response["message"]) {
			this.message = response["message"]
		} else {
			this.message = codeToJSON(response["code"])
		}

		switch (this.message) {
			case "Failed to fetch":
				this.message = "网络连接失败";
				break;
		}

		this.code = response["code"];
	}
}

function codeToJSON(responseElement: number): string {
	switch (responseElement) {
		default:
			return "非预期错误"
	}
}