# Gunakan base image golang
FROM golang:1.20

# Set working directory di dalam container
WORKDIR /app

# Copy go.mod dan go.sum untuk instalasi dependencies
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy seluruh kode aplikasi ke dalam container
COPY . .

# Compile aplikasi
RUN go build -o main .

# Ekspos port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
