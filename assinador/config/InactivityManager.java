package br.ufg.inf.assinador.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;

@Component
public class InactivityManager {

    private LocalDateTime lastActivityTime = LocalDateTime.now();
    private final ApplicationContext context;
    private final int timeoutMinutes;

    public InactivityManager(ApplicationContext context, 
                             @Value("${assinador.timeout.minutes:0}") int timeoutMinutes) {
        this.context = context;
        this.timeoutMinutes = timeoutMinutes;
    }

  
    public void registerActivity() {
        this.lastActivityTime = LocalDateTime.now();
    }

    @Scheduled(fixedDelay = 60000)
    public void checkInactivity() {
        if (timeoutMinutes <= 0) {
            return;
        }

        long minutesInactive = ChronoUnit.MINUTES.between(lastActivityTime, LocalDateTime.now());
        
        if (minutesInactive >= timeoutMinutes) {
            System.out.println("⏳ Servidor ocioso por " + timeoutMinutes + " minutos. Encerrando por inatividade...");
            SpringApplication.exit(context, () -> 0);
            System.exit(0);
        }
    }
}