package com.resiliency;

import breaker.CircuitBreaker;

public class Main {
    public static void main(String[] args) throws InterruptedException {
        int failureThreshold = 3;
        int successThreshold = 2;
        long intervalThreshold = 2000; // 2 detik

        CircuitBreaker breaker = new CircuitBreaker(failureThreshold, successThreshold, intervalThreshold);

        System.out.println("=== MULAI SIMULASI CIRCUIT BREAKER ===");

        // Simulasi 3 kegagalan berturut-turut
        for (int i = 1; i <= 3; i++) {
            System.out.println("[FAILURE] ke-" + i);
            breaker.onFailure();
            System.out.println("State: " + breaker.getCurrentState());
        }

        // Setelah mencapai ambang kegagalan, circuit breaker masuk ke state OPEN
        System.out.println("\n[INFO] Mencoba akses saat state OPEN...");
        boolean permission = breaker.acquirePermission();
        System.out.println("acquirePermission(): " + permission);
        System.out.println("State: " + breaker.getCurrentState());

        // Tunggu hingga intervalThreshold habis
        System.out.println("\n[INFO] Menunggu " + intervalThreshold + "ms agar bisa masuk HALF_OPEN...");
        Thread.sleep(intervalThreshold + 100);

        System.out.println("[INFO] Mencoba akses lagi setelah interval...");
        permission = breaker.acquirePermission();
        System.out.println("acquirePermission(): " + permission);
        System.out.println("State: " + breaker.getCurrentState());

        // Simulasi 2 keberhasilan berturut-turut di HALF_OPEN
        for (int i = 1; i <= 2; i++) {
            System.out.println("[SUCCESS] ke-" + i);
            breaker.onSuccess();
            System.out.println("State: " + breaker.getCurrentState());
        }

        // Circuit breaker harus kembali ke CLOSE
        System.out.println("\n[INFO] Status akhir:");
        System.out.println("State: " + breaker.getCurrentState());
        System.out.println("=== SELESAI ===");
    }
}