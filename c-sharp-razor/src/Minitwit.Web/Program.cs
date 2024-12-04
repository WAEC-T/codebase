using FluentValidation;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Infrastructure;
using Minitwit.Infrastructure.Repository;
using Minitwit.Razor.Pages;
using Minitwit.Web;

public class Program
{
    public static void Main(string[] args)
    {
        var builder = WebApplication.CreateBuilder(args);

        // Add services to the container.
        builder.Services.AddRazorPages(options =>
        {
            options.Conventions.AddAreaPageRoute("Identity", "/Account/Register", "/register");
            options.Conventions.AddAreaPageRoute("Identity", "/Account/Login", "/login");
            options.Conventions.AddAreaPageRoute("Identity", "/Account/Logout", "/logout");
        });
        
        builder.Services.AddMvc(options =>
        {
            options.Filters.Add(new IgnoreAntiforgeryTokenAttribute());
        });
        
        ProgramOptions.AddProgramOptions(builder);
        ProgramOptions.AddIdendity(builder);
        ProgramOptions.AddDatabase(builder);

        //API Controllers
        builder
            .Services.AddControllers()
            .AddJsonOptions(options =>
            {
                options.JsonSerializerOptions.PropertyNameCaseInsensitive = true;
                options.JsonSerializerOptions.PropertyNamingPolicy = null;
                options.JsonSerializerOptions.IgnoreNullValues = true;
            });
        
        builder.Services.ConfigureApplicationCookie(options =>
        {
            options.Cookie.SecurePolicy = CookieSecurePolicy.None; // Allow cookies over HTTP
            options.Cookie.HttpOnly = true;
            options.Cookie.SameSite = SameSiteMode.Lax; // Use Lax or None depending on requirements
            options.ExpireTimeSpan = TimeSpan.FromMinutes(60); // Extend expiration for testing
            options.SlidingExpiration = true; // Enable sliding expiration
        });

        builder.Services.AddSession(options =>
        {
            options.Cookie.Name = ".Minitwit.Web.Session";
            options.IdleTimeout = TimeSpan.FromMinutes(30);
            options.Cookie.HttpOnly = true;
            options.Cookie.IsEssential = true;
            options.Cookie.SecurePolicy = CookieSecurePolicy.None; // Allow cookies over HTTP
            options.Cookie.SameSite = SameSiteMode.Lax;
        });

        // Dependency Injection
        builder.Services.AddScoped<IAuthorRepository, AuthorRepository>();
        builder.Services.AddScoped<IValidator<CreateMessage>, MessageCreateValidator>();
        builder.Services.AddScoped<IMessageRepository, MessageRepository>();
        builder.Services.AddScoped<IMessageService, MinitwitService>();
        builder.Services.AddScoped<IFollowRepository, FollowRepository>();
        builder.Services.AddScoped<ILatestRepository, LatestRepository>();

        builder.Logging.ClearProviders();
        builder.Logging.AddConsole();
        builder.Logging.AddDebug();

        builder.Services.AddSession(options =>
        {
            options.Cookie.Name = ".Minitwit.Web.Session";
            options.IdleTimeout = TimeSpan.FromMinutes(10);
            options.Cookie.HttpOnly = false;
            options.Cookie.IsEssential = true;
        });

        var app = builder.Build();

        // Apply migrations
        using (var scope = app.Services.CreateScope())
        {
            var services = scope.ServiceProvider;

            try
            {   
                var context = services.GetRequiredService<MinitwitDbContext>();
                context.Database.Migrate();
            }
            catch (Exception e)
            {
                Console.WriteLine($"Error applying migrations: {e.Message}");
            }
        }

        if (!app.Environment.IsDevelopment())
        {
            app.UseExceptionHandler("/Error");
            app.UseHsts();
        }
        app.UseStaticFiles();
        
        app.UseRouting();

        app.UseAuthentication();
        app.UseAuthorization();
        
        app.UseSession(); 
        app.MapControllers();
        app.MapRazorPages();
        
        app.Run();
    }
}
