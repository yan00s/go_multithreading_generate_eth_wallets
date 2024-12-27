## Ethereum Wallet Generation Benchmark with Threading

This project is a multi-threaded Ethereum wallet generator, which can also be used as a cpu benchmark for testing different hardware setups.

### Benchmark Results:

#### 2024-12-25 on **Snapdragon 720G** (Golang)
- **Wallets Created:** 100,000  
- **Threads:** 8  
- **Time:** 125.57 seconds  

#### 2024-12-25 on **RK3588** (Golang)
- **Wallets Created:** 100,000  
- **Threads:** 8  
- **Time:** 78.230 seconds  

#### 2024-12-22 on **Intel i5-6300HQ** (Golang)
- **Wallets Created:** 100,000  
- **Threads:** 4  
- **Time:** 59.587 seconds  

#### Pre-2024 on **AMD Ryzen 5 4600H**
- **Wallets Created:** 100,000  
- **Threads:** 12  
- **Time:** 19.893 seconds  
