# Use a Node.js image (Node 20 LTS to avoid localStorage issues)
FROM node:20

# Set the working directory
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy everything to the container
COPY . .

# Expose port
EXPOSE 3000

# Start the development server
CMD ["npm", "start"]