package dk.itu.minitwit.controller;

// import dk.itu.minitwit.database.SQLite;
import dk.itu.minitwit.database.PostgreSQL;
import dk.itu.minitwit.domain.AddMessage;
import dk.itu.minitwit.domain.Login;
import dk.itu.minitwit.domain.Register;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpSession;
import org.apache.commons.codec.digest.DigestUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import java.sql.SQLException;
import java.sql.Timestamp;
import java.text.SimpleDateFormat;
import java.util.*;

@Controller
public class MiniTwitController {

    @Autowired
    // protected SQLite sqLite;
    protected PostgreSQL postgreSQL;

    @Autowired
    PasswordEncoder passwordEncoder;

    final static int PER_PAGE = 30;

    SimpleDateFormat sdf = new SimpleDateFormat("MMM dd,yyyy HH:mm:ss");
    Logger logger = LoggerFactory.getLogger(MiniTwitController.class);


    @GetMapping("/")
    public String timeline(Model model, HttpServletRequest request) {
        
        try {
            logger.info("Request ID: %s -- Received %s request on path: '%s'"
                    .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));
            HttpSession session = request.getSession(false);
            boolean loggedIn = addUserToModel(model, session);
            if (!loggedIn) {
                logger.info("Request ID: %s -- User not logged in, redirecting to '/public'".formatted( model.getAttribute("requestID")));
                return "redirect:/public";
            }

            List<Object> args = new ArrayList<>();

            args.add(session.getAttribute("user_id"));
            args.add(session.getAttribute("user_id"));
            args.add(PER_PAGE);
            List<Map<String, Object>> messages = null;
            long before = System.currentTimeMillis();
            logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
            messages = postgreSQL.queryDb(
                    "select messages.*, " +
                            "users.* from messages, " +
                            "users where messages.flagged = 0 " +
                            "and messages.author_id = users.user_id " +
                            "and (users.user_id = ? or users.user_id in (select whom_id from followers where who_id = ?))" +
                            "order by messages.pub_date desc limit ?"
                    , args);

            long after = System.currentTimeMillis();
            logger.info("Request ID: %s -- Queried database in %.2f seconds"
                    .formatted( model.getAttribute("requestID"), getDuration(before, after)));
            addDatesAndGravatarURLs(messages);

            model.addAttribute("my", "true");
            model.addAttribute("messages", messages);
            model.addAttribute("messagesSize", messages.size());
        } catch (SQLException | NullPointerException e) {
            logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                    "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
        }
        String template = "timeline.html";
        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        return template;
    }

    @RequestMapping(value = "/public", method = RequestMethod.GET)
    public Object publicTimeline(Model model, HttpServletRequest request) {
        try {
            logger.info("Request ID: %s -- Received %s request on path: '%s'"
                    .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));
            model.addAttribute("public", "true");
            HttpSession session = request.getSession(false);
            addUserToModel(model, session);

            List<Object> args = new ArrayList<>();
            args.add(PER_PAGE);
            List<Map<String, Object>> messages = null;
            long before = System.currentTimeMillis();
            logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));

            messages = postgreSQL.queryDb(
                    "select messages.*, users.* from messages, users " +
                            // "where message.flagged = 0 and message.author_id = user.user_id " +
                            "where messages.author_id = users.user_id " +
                            "order by messages.pub_date desc limit ?", args);

            long after = System.currentTimeMillis();
            logger.info("Request ID: %s -- Queried database in %.2f seconds"
                    .formatted( model.getAttribute("requestID"), getDuration(before, after)));
            addDatesAndGravatarURLs(messages);

            model.addAttribute("messages", messages);
            model.addAttribute("messagesSize", messages.size());
        } catch (SQLException | NullPointerException e){
            logger.error(
                    "Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                            "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
        }
        String template = "timeline.html";
        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        return template;
    }


    @GetMapping("/user/{username}")
    public String userTimeLine(@PathVariable("username") String username, HttpServletRequest request, Model model) {
        logger.info("Request ID: %s -- Received %s request on path: '%s'"
                .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));
        HttpSession session = request.getSession(false);
        model.addAttribute("public", "false");
        model.addAttribute("username", username);
        boolean loggedIn = addUserToModel(model, session);

        List<Object> arg = new ArrayList<>();
        arg.add(username);
        List<Map<String, Object>> users = null;
        long before = System.currentTimeMillis();
        logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
        try {
            users = postgreSQL.queryDb("select * from users where users.username = ?", arg);
        } catch (SQLException e) {
            logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                    "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
        }
        long after = System.currentTimeMillis();
        logger.info("Request ID: %s -- Queried database in %.2f seconds"
                .formatted( model.getAttribute("requestID"), getDuration(before, after)));

        if (users.size() == 0) {
            throw new ResponseStatusException(
                    HttpStatus.NOT_FOUND, "User not found"
            );
        }

        if (loggedIn) {
            int otherId = 0;
            before = System.currentTimeMillis();
            logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
            try {

                otherId = getUserID(username);
            } catch (SQLException e) {
                logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                        "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
            }
            after = System.currentTimeMillis();
            logger.info("Request ID: %s -- Queried database in %.2f seconds"
                    .formatted( model.getAttribute("requestID"), getDuration(before, after)));

            if ((int) session.getAttribute("user_id") == otherId) {
                model.addAttribute("self", "true");
                model.addAttribute("followed", "false");
            } else {
                List<Object> args = new ArrayList<>();
                args.add(session.getAttribute("user_id"));
                args.add(otherId);
                List<Map<String, Object>> followed = null;
                before = System.currentTimeMillis();
                logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
                try {
                    followed = postgreSQL.queryDb("select * from followers where followers.who_id = ? and followers.whom_id = ?", args);
                } catch (SQLException e) {
                    logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                            "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
                }
                after = System.currentTimeMillis();
                logger.info("Request ID: %s -- Queried database in %.2f seconds"
                        .formatted( model.getAttribute("requestID"), getDuration(before, after)));
                model.addAttribute("followed", followed.size() > 0 ? "true" : "false");
            }
        }
        List<Object> args = new ArrayList<>();
        args.add(username);
        args.add(PER_PAGE);
        try {
        List<Map<String, Object>> messages = null;
        before = System.currentTimeMillis();
        logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));

            messages = postgreSQL.queryDb(
                    "select messages.*, " +
                            "users.* from messages, users " +
                            "where users.username = ? " +
                            "and messages.flagged = 0 " +
                            "and messages.author_id = users.user_id " +
                            "order by messages.pub_date desc limit ?"
                    , args);

        after = System.currentTimeMillis();
        logger.info("Request ID: %s -- Queried database in %.2f seconds"
                .formatted( model.getAttribute("requestID"), getDuration(before, after)));
        addDatesAndGravatarURLs(messages);

        model.addAttribute("messages", messages);
        model.addAttribute("messagesSize", messages.size());
        } catch (SQLException | NullPointerException e) {
            logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                    "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
        }
        String template = "timeline.html";
        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        return template;
    }


    @GetMapping("/{username}/follow")
    public String followUser(@PathVariable("username") String username, HttpServletRequest request, Model model, RedirectAttributes redirectAttributes) {
        logger.info("Request ID: %s -- Received %s request on path: '%s'"
                .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));
        HttpSession session = request.getSession(false);
        boolean loggedIn = addUserToModel(model, session);

        if (!loggedIn) {
            logger.info("Request ID: %s -- User not logged in, redirecting to '/login'");
            return "redirect:/login";
        }

        Integer whomId;
        try {
            whomId = getUserID(username);
        } catch (SQLException e) {
            logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                    "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
            return "";
        }
        List<Object> args = new ArrayList<>();
        args.add(session.getAttribute("user_id"));
        args.add(whomId);
        long before = System.currentTimeMillis();
        logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
        try {
            postgreSQL.updateDb("insert into followers (who_id, whom_id) values (?, ?)", args);
        } catch (SQLException e) {
            logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() + "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
        }
        long after = System.currentTimeMillis();
        logger.info("Request ID: %s -- Queried database in %.2f seconds"
                .formatted( model.getAttribute("requestID"), getDuration(before, after)));
        String template = "redirect:/user/" + username;
        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        redirectAttributes.addFlashAttribute("flashMessage", "You are now following " + username);
        return template;
    }

    @GetMapping("/{username}/unfollow")
    public String unfollowUser(@PathVariable("username") String username, HttpServletRequest request, Model model, RedirectAttributes redirectAttributes) {
        logger.info("Request ID: %s -- Received %s request on path: '%s'"
                .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));

        HttpSession session = request.getSession(false);
        boolean loggedIn = addUserToModel(model, session);

        if (!loggedIn) {
            logger.info("Request ID: %s -- User not logged in, redirecting to '/login'");
            return "redirect:/login";
        }
        Integer whomId;
        try {
            whomId = getUserID(username);
        } catch (SQLException e) {
            logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                    "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
            return "";
        }
        if (whomId == null) {
            logger.info("Request ID: %s -- User: %s not found - cannot follow - redirecting to public"
                    .formatted( model.getAttribute("requestID"), username));
            return "redirect:/public";
        }


        List<Object> args = new ArrayList<>();
        args.add(session.getAttribute("user_id"));
        args.add(whomId);
        long before = System.currentTimeMillis();
        logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
        try {
            postgreSQL.updateDb("delete from followers where who_id=? and whom_id=?", args);
        } catch (SQLException e) {
            logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() + "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
        }
        long after = System.currentTimeMillis();
        logger.info("Request ID: %s -- Queried database in %.2f seconds"
                .formatted( model.getAttribute("requestID"), getDuration(before, after)));
        redirectAttributes.addFlashAttribute("flashMessage", "You are no longer following " + username);
        String template = "redirect:/user/" + username;
        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        return template;
    }

    @PostMapping("/add_message")
    public String addMessage(AddMessage text, HttpServletRequest request, Model model, RedirectAttributes redirectAttributes) {
        logger.info("Request ID: %s -- Received %s request on path: '%s'"
                .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));
        HttpSession session = request.getSession(false);
        boolean loggedIn = addUserToModel(model, session);
        if (!loggedIn) {
            logger.info("Request ID: %s -- User not logged in, redirecting to '/login'");
            return "redirect:/login";
        }

        if (text.getText() != null && !text.getText().isEmpty()) {
            List<Object> args = new ArrayList<>();
            args.add(session.getAttribute("user_id"));
            args.add(text.getText());
            // Timestamp timestamp = new Timestamp(System.currentTimeMillis() / 1000);
            args.add(new Timestamp(System.currentTimeMillis() / 1000));
            long before = System.currentTimeMillis();
            logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
            try {
                postgreSQL.updateDb("insert into messages (author_id, text, pub_date, flagged) values (?, ?, ?, 0)", args);
            } catch (SQLException e) {
                logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                        "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
            }
            long after = System.currentTimeMillis();
            logger.info("Request ID: %s -- Queried database in %.2f seconds"
                    .formatted( model.getAttribute("requestID"), getDuration(before, after)));
            redirectAttributes.addFlashAttribute("flashMessage", "Your message was recorded");
        } else {
            redirectAttributes.addFlashAttribute("error", "Message cannot be empty!");
            model.addAttribute("error", "Message cannot be empty!");
        }

        String referer = request.getHeader("Referer");
        String template = "";

        if (referer != null){
            if (referer.contains("/public"))
                template = "redirect:/public";
            else
                template = "redirect:/";
        } else {
            template = "redirect:/public";
        }

        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        return template;
    }


    @RequestMapping(value = "/login", method = {RequestMethod.GET, RequestMethod.POST})
    public String login(@ModelAttribute Login login, Model model, HttpServletRequest request, RedirectAttributes redirectAttributes) {
        logger.info("Request ID: %s -- Received %s request on path: '%s'"
                .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));


        if ("POST".equals(request.getMethod())) {
            //query db med user/pass fra login objekt
            List<Object> args = new ArrayList<>();
            args.add(login.getUsername());
            List<Map<String, Object>> s = null;
            long before = System.currentTimeMillis();
            logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
            try {
                s = postgreSQL.queryDb("select * from users where username = ?", args);
            } catch (SQLException e) {
                logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                        "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
            }
            long after = System.currentTimeMillis();
            logger.info("Request ID: %s -- Queried database in %.2f seconds"
                    .formatted( model.getAttribute("requestID"), getDuration(before, after)));
            if (s.isEmpty()) { // Wrong credentials, return 404
                logger.info("Request ID: %s -- User not logged in, invalid credentials");
                model.addAttribute("error", "Invalid credentials");
                return "login.html";
            } else if (!passwordEncoder.matches(login.getPassword(), (String) s.get(0).get("pw_hash"))) {
                logger.info("Request ID: %s -- User not logged in, invalid credentials");
                model.addAttribute("error", "Invalid credentials");
                return "login.html";
            } else { // change redirects to my timeline
                // Session
                logger.info("Request ID: %s -- User logged in, redirecting to '/'");
                
                request.getSession().setAttribute("user", login.getUsername());
                request.getSession().setAttribute("user_id", s.get(0).get("user_id"));

                redirectAttributes.addFlashAttribute("flashMessage", "You were logged in");
                String template = "redirect:/";
                logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
                return template;
            }
        }
        if (request.getParameter("logout") != null) {
            redirectAttributes.addFlashAttribute("flashMessage", "You were logged out");
            return "redirect:/public";
        }
        String template = "login.html";
        logger.info("Testtest: Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        return template;
    }


    @RequestMapping(value = "/register", method = {RequestMethod.GET, RequestMethod.POST})
    public String register(@ModelAttribute Register register, Model model, HttpServletRequest request, RedirectAttributes redirectAttributes) {
        HttpSession session = request.getSession(false);
        addUserToModel(model, session);
        Object user = model.getAttribute("user");

        if (model.getAttribute("user") != null && model.getAttribute("user") != "" && model.getAttribute("user") != "false") {
            logger.info("Request ID: %s -- User already logged in, redirecting to '/public'".formatted( model.getAttribute("requestID")));
            return "redirect:/";
        }

        logger.info("Request ID: %s -- Received %s request on path: '%s'"
                .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));

        if ("POST".equals(request.getMethod())) {
            if ("".equals(register.getUsername())) {
                logger.info("Request ID: %s -- User not registered - no username in input");
                model.addAttribute("error", "You have to enter a username");
                return "register.html";
            } else if ("".equals(register.getEmail()) || !register.getEmail().matches("^[\\w.-]+@[\\w.-]+\\.[a-zA-Z]{2,}$")) {
                logger.info("Request ID: %s -- User not registered - invalid email address entered");
                model.addAttribute("error", "You have to enter a valid email address");
                return "register.html";
            } else if ("".equals(register.getPassword())) {
                logger.info("Request ID: %s -- User not registered - no password in input");
                model.addAttribute("error", "You have to enter a password");
                return "register.html";
            } else if (!register.getPassword2().equals(register.getPassword())) {
                logger.info("Request ID: %s -- User not registered - the two passwords do not match");
                model.addAttribute("error", "The two passwords do not match");
                return "register.html";
            } else {
                List<Object> args = new ArrayList<>();
                args.add(register.getUsername());
                args.add(register.getEmail());
                args.add(passwordEncoder.encode(register.getPassword()));
                
                long before = System.currentTimeMillis();
                logger.info("Request ID: %s -- Querying database...".formatted( model.getAttribute("requestID")));
                try {
                    postgreSQL.updateDb("insert into users (username, email, pw_hash) values (?, ?, ?)", args);
                } catch (SQLException e) {
                    logger.error("Request ID: %s -- Encountered error while querying database: " + e.getMessage() +
                    "\n" + Arrays.toString(e.getStackTrace()).formatted( model.getAttribute("requestID")));
                }
                long after = System.currentTimeMillis();
                logger.info("Request ID: %s -- Queried database in %.2f seconds"
                .formatted( model.getAttribute("requestID"), getDuration(before, after)));
                
                String flashMessage = "You were successfully registered and can login now"; 
                //model.addAttribute("flashMessage", flashMessage);
                redirectAttributes.addFlashAttribute("flashMessage", flashMessage);
                String template = "redirect:/login";
                logger.info("Request ID: %s -- User registered - Returning template: %s".formatted( model.getAttribute("requestID"), template));
                return template;
            }
        }

        String template = "register.html";
        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        return template;

    }

    @GetMapping("/logout")
    public String logout(HttpServletRequest request, Model model, RedirectAttributes redirectAttributes) {
        logger.info("Request ID: %s -- Received %s request on path: '%s'"
                .formatted( model.getAttribute("requestID"), request.getMethod(), request.getRequestURI()));
                
        request.getSession().invalidate();
        addUserToModel(model, request.getSession(false));
        logger.info("Request ID: %s -- Session invalidated".formatted( model.getAttribute("requestID")));
        String template = "redirect:/public";
        logger.info("Request ID: %s -- Returning template: %s".formatted( model.getAttribute("requestID"), template));
        redirectAttributes.addFlashAttribute("flashMessage", "You were logged out");
        return template;
    }

    @ModelAttribute("requestID")
    public String requestId() {
        return UUID.randomUUID().toString();
    }

    public int getUserID(String username) throws SQLException {
        List<Object> args = new ArrayList<>();
        args.add(username);
        List<Map<String, Object>> userIDs;
        userIDs = postgreSQL.queryDb("select user_id from users where username = ?", args);
        return ((int) userIDs.get(0).get("user_id"));
    }

    public static boolean addUserToModel(Model model, HttpSession session) {
        Object user = null;
        try {

            user = session.getAttribute("user");
        } catch (Exception e){
            model.addAttribute("user", "false");
            System.out.println("User not logged in");
            return false;
        }
        if (user == "" || user == null)
        {
            model.addAttribute("user", "false");
            System.out.println("User not logged in");
            return false;
        }
        if (session != null) {
            model.addAttribute("user", session.getAttribute("user"));
            return true;
        } else {
            model.addAttribute("user", "false");
            return false;
        }
    }

    public void addDatesAndGravatarURLs(List<Map<String, Object>> messages) {
        messages.forEach(obj -> {
            String email = (String) obj.get("email");
            Timestamp timestamp = (Timestamp) obj.get("pub_date");
            Long created = timestamp.getTime();
            obj.put("gravatar_url", "https://www.gravatar.com/avatar/" + getMD5Hash(email.toLowerCase().strip()) + "?d=identicon&s=80");
            Date d = new Date((created) * 1000);
            obj.put("date_time", sdf.format(d));
        });
    }

    public String getMD5Hash(String email) {
        return DigestUtils.md5Hex(email.toLowerCase());
    }

    public double getDuration(long before, long after) {
        return ((double) after - (double) before) / 1000;
    }
}
