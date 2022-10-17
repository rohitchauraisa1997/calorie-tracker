# Use a Node 16 base image
FROM node:16-alpine 
COPY . .
RUN rm -rf .git/hooks
RUN npm install 
EXPOSE 3000
CMD ["npm", "start"]