package breaker;

import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.Setter;
import lombok.AccessLevel;
import lombok.extern.slf4j.Slf4j;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Slf4j
@Setter
@Getter
@AllArgsConstructor
@RequiredArgsConstructor
public class CircuitBreaker {
    Logger logger = LoggerFactory.getLogger(CircuitBreaker.class);

    public enum State {
        CLOSE, OPEN, HALF_OPEN
    }

    private final Integer failureThreshold;
    private final Integer successThreshold;
    private final long intervalThreshold;

    @Setter(AccessLevel.PRIVATE)
    private State currentState = State.CLOSE;

    private final AtomicInteger failureCount = new AtomicInteger(0);
    private final AtomicInteger successCount = new AtomicInteger(0);
    private final AtomicLong lastFailureTimestamp = new AtomicLong(0);

    public synchronized boolean acquirePermission() {
        try {
            switch (currentState) {
                case CLOSE, HALF_OPEN:
                    return true;
                case OPEN:
                    return this.handleOpenState();
                default:
                    throw new Exception("state undefined");
            }
        } catch (Exception e) {
            e.printStackTrace();
            return false;
        }
    }

    private synchronized boolean handleOpenState() {
        long diffTime = System.currentTimeMillis() - this.getLastFailureTimestamp().get();
        if (diffTime >= this.getIntervalThreshold()) {
            this.setCurrentState(State.HALF_OPEN);
            this.reset();

            return true;
        }

        return false;
    }

    public synchronized void onSuccess() {
        if (this.getCurrentState() == State.HALF_OPEN) {
            boolean isSuccess = (this.getSuccessCount().incrementAndGet() >= this.getSuccessThreshold());
            if (isSuccess) {
                this.setCurrentState(State.CLOSE);
                log.debug("CircuitBreaker: state recovered to Closed");
                this.reset();
            }

            if (this.getFailureCount().incrementAndGet() >= this.getFailureThreshold()) {
                this.setCurrentState(State.OPEN);
                this.reset();
                this.getLastFailureTimestamp().set(System.currentTimeMillis());
            }
        }

        if (System.currentTimeMillis() - this.getLastFailureTimestamp().get() > this.getIntervalThreshold()) {
            this.getFailureCount().set(0);
        }
    }

    public synchronized void onFailure() {
        currentState = this.getCurrentState();
        if (currentState == State.CLOSE || currentState == State.HALF_OPEN) {
            int failureCount = this.getFailureCount().incrementAndGet();
            if (failureCount >= this.getFailureThreshold()) {
                this.setCurrentState(State.OPEN);
                this.reset();
                this.getLastFailureTimestamp().set(System.currentTimeMillis());
            }
        }
    }

    private synchronized void reset() {
        this.getSuccessCount().set(0);
        this.getFailureCount().set(0);
        this.getLastFailureTimestamp().set(0);
    }
}
