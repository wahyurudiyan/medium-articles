package breaker;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("Unit Test for CircuitBreaker")
class CircuitBreakerTest {

    private CircuitBreaker breaker;

    private final int failureThreshold = 3;
    private final int successThreshold = 2;
    private final long intervalThreshold = 1000; // 1 detik

    @BeforeEach
    void setup() {
        breaker = new CircuitBreaker(failureThreshold, successThreshold, intervalThreshold);
    }

    @Test
    @DisplayName("State awal seharusnya CLOSE")
    void testInitialStateIsClose() {
        assertEquals(CircuitBreaker.State.CLOSE, breaker.getCurrentState());
    }

    @Test
    @DisplayName("acquirePermission() harus true saat state CLOSE")
    void testAcquirePermissionInCloseState() {
        assertTrue(breaker.acquirePermission());
    }

    @Test
    @DisplayName("State berubah menjadi OPEN setelah mencapai ambang batas gagal")
    void testTransitionToOpenAfterFailureThreshold() {
        breaker.onFailure(); // 1
        breaker.onFailure(); // 2
        breaker.onFailure(); // 3 -> should open

        assertEquals(CircuitBreaker.State.OPEN, breaker.getCurrentState());
    }

    @Test
    @DisplayName("acquirePermission() harus false saat state OPEN sebelum timeout")
    void testAcquirePermissionInOpenBeforeTimeout() {
        breaker.onFailure();
        breaker.onFailure();
        breaker.onFailure();

        assertFalse(breaker.acquirePermission());
    }

    @Test
    @DisplayName("State berubah ke HALF_OPEN setelah waktu interval habis")
    void testTransitionToHalfOpenAfterInterval() throws InterruptedException {
        breaker.onFailure();
        breaker.onFailure();
        breaker.onFailure();

        Thread.sleep(intervalThreshold + 100); // tunggu timeout

        assertTrue(breaker.acquirePermission());
        assertEquals(CircuitBreaker.State.HALF_OPEN, breaker.getCurrentState());
    }

    @Test
    @DisplayName("State kembali ke CLOSE dari HALF_OPEN setelah sukses beruntun")
    void testTransitionToCloseFromHalfOpenOnSuccesses() throws InterruptedException {
        breaker.onFailure();
        breaker.onFailure();
        breaker.onFailure();

        Thread.sleep(intervalThreshold + 50);
        breaker.acquirePermission(); // Move to HALF_OPEN

        breaker.onSuccess(); // 1
        breaker.onSuccess(); // 2 -> should close

        assertEquals(CircuitBreaker.State.CLOSE, breaker.getCurrentState());
    }

    @Test
    @DisplayName("Counter failure harus direset setelah interval")
    void testResetFailureCountAfterInterval() throws InterruptedException {
        breaker.onFailure(); // failureCount = 1

        Thread.sleep(intervalThreshold + 50);
        breaker.onSuccess(); // harus reset failureCount

        assertEquals(0, breaker.getFailureCount().get());
    }

    @Test
    @DisplayName("State kembali ke OPEN jika gagal saat HALF_OPEN")
    void testTransitionBackToOpenFromHalfOpenOnFailure() throws InterruptedException {
        breaker.onFailure();
        breaker.onFailure();
        breaker.onFailure();

        Thread.sleep(intervalThreshold + 50);
        breaker.acquirePermission(); // HALF_OPEN

        breaker.onFailure();
        breaker.onFailure();
        breaker.onFailure(); // Failure saat HALF_OPEN -> back to OPEN

        assertEquals(CircuitBreaker.State.OPEN, breaker.getCurrentState());
    }
}