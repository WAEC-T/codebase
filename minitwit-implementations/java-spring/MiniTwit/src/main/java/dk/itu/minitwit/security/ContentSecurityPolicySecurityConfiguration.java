package dk.itu.minitwit.security;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;

@EnableWebSecurity
@Configuration
public class ContentSecurityPolicySecurityConfiguration {
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
            .csrf().disable()
            .authorizeHttpRequests((authz) -> authz
                .anyRequest().permitAll()
            )
            .headers()
            .xssProtection()
            .and()
            .contentSecurityPolicy("form-action 'self' always; default-src 'none'; script-src 'self'; style-src 'self' *.academicweapons.dk/*; font-src 'self'; connect-src 'self'; img-src 'self' *.gravatar.com; frame-src 'none'; frame-ancestors 'none'; media-src 'none'; object-src 'none'; manifest-src 'none'; worker-src 'none';");

        return http.build();
    }
    
}

