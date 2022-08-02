import { defineEndpoint } from '@directus/extensions-sdk';
import { randomInt } from 'crypto';

export default defineEndpoint((router, { services, exceptions }) => {
	const { ItemsService } = services;
	const { ServiceUnavailableException } = exceptions;
	router.get('/:code', (req, res, next) => {
		const urlService = new ItemsService('urls',{schema: req.schema});
		urlService
			.readByQuery({ sort: ['url'],filter:{code: {_eq:req.params.code},active: {_eq:'true'}}, fields: ['*'] })
			.then((results:Array<any>) => {
				if(results.length==0){
					res.redirect("https://www.youtube.com/watch?v=dQw4w9WgXcQ");
					return;
				}
				const urls = [];
				for(let res of results){
					for(let i=0;i<res.likelyhood;i++){
						urls.push(res.url)
					}
					
				}

				res.redirect(urls[randomInt(urls.length)])
			})
			.catch((error:any) => {
				return next(new ServiceUnavailableException(error.message));
			});
	});

});
