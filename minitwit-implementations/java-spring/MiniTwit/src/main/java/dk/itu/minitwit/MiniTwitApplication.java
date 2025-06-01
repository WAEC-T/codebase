package dk.itu.minitwit;


import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;




@SpringBootApplication(scanBasePackages={"dk.itu.minitwit.database","dk.itu.minitwit.controller","dk.itu.minitwit.security"} ,exclude = {
        org.springframework.boot.autoconfigure.security.servlet.SecurityAutoConfiguration.class,
        org.springframework.boot.actuate.autoconfigure.security.servlet.ManagementWebSecurityAutoConfiguration.class}
)
public class MiniTwitApplication {


    public static void main(final String[] args) {
        // Class.forName("org.postgresql.Driver");
        
        
        SpringApplication.run(MiniTwitApplication.class, args);
    }

    @Bean
    public BCryptPasswordEncoder bCryptPasswordEncoder() {
        return new BCryptPasswordEncoder();
    }

}
