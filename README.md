# grpc-realtime-chat - A chat application using gRPC and NATS

This is an example project using gRPC and NATS to build one realtime chat.
It will contain gRPC services in go and one client using react native.

# Folder structure

- `/proto`: Contain all protocol buffers files
- `/server/cmd`: Contain all executables for services and the load balancer
- `/mobile`: (TODO) The clone of [WhatsApp in React Native](https://github.com/VctrySam/whatsapp) implemented using the gRPC client
- `/electron-app`: The clone of [Chatron - Chat application with react and electron](https://github.com/LennyBoyatzis/Chatron). One electron chat app using the grpc node client.
