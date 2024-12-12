import './App.css';
import { binary_to_base58 } from "base58-js";

function App() {
  const getProvider = () => {
    if ('phantom' in window) {
      const provider = window.phantom?.solana;
  
      if (provider?.isPhantom) {
        return provider;
      }
    }
  };

  const handleLogin = async () => {
    const provider = getProvider();
    try {
        const resp = await provider.connect();
        const publicAddress = resp.publicKey.toString();
        
        const message = "Sign me";
        const encodedMessage = new TextEncoder().encode(message);
        const signature = await provider.signMessage(encodedMessage, "hex");

        await handleAuthenticate(publicAddress, binary_to_base58(signature.signature))

    } catch (error) {
    }
  }

    const handleAuthenticate = async (publicAddress, signature) => {
      console.log("publicAddress: ", publicAddress);
      console.log("signature: ", signature);
  
      // 要发送到 Go API 的请求体
      const requestBody = {
          publicAddress: publicAddress,
          signature: signature,
      };
  
      try {
          // 发送 POST 请求到 Go API
          const response = await fetch('http://localhost:8080/api/verify', {
              method: 'POST',
              headers: {
                  'Content-Type': 'application/json',
              },
              body: JSON.stringify(requestBody),
          });
  
          // 检查响应状态
          if (!response.ok) {
              throw new Error(`HTTP error! status: ${response.status}`);
          }
  
          // 解析响应 JSON
          const data = await response.json();
          
          // 处理返回的 JWT
          if (data.token) {
              console.log("JWT Token:", data.token);
              // 这里可以将 JWT 存储在本地存储或状态管理中
              localStorage.setItem('jwtToken', data.token);
          } else {
              console.error("No token returned from the server.");
          }
      } catch (error) {
          console.error("Error during authentication:", error);
      }
    }
    
  return (
    <div className="App">
      <button onClick={handleLogin}>Login - solana</button>
    </div>
  );
}

export default App;
