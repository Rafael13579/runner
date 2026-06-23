package br.ufg.inf.assinador.config;

import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class ActivityFilter implements Filter {

    private final InactivityManager inactivityManager;

    public ActivityFilter(InactivityManager inactivityManager) {
        this.inactivityManager = inactivityManager;
    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) 
            throws IOException, ServletException {
        
        inactivityManager.registerActivity(); 
        chain.doFilter(request, response);  
    }
}