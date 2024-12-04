// Licensed to the .NET Foundation under one or more agreements.
// The .NET Foundation licenses this file to you under the MIT license.
#nullable disable

using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Text.Encodings.Web;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.UI.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Extensions.Logging;
using Minitwit.Core.Entities;

namespace Minitwit.Web.Areas.Identity.Pages.Account
{   
    [IgnoreAntiforgeryToken]
    public class RegisterModel : PageModel
    {
        private readonly SignInManager<Author> _signInManager;
        private readonly UserManager<Author> _userManager;
        private readonly IUserStore<Author> _userStore;
        private readonly IUserEmailStore<Author> _emailStore;
        private readonly ILogger<RegisterModel> _logger;
        private readonly IEmailSender _emailSender;

        public RegisterModel(
            UserManager<Author> userManager,
            SignInManager<Author> signInManager,
            ILogger<RegisterModel> logger)
        {
            _userManager = userManager;
            _signInManager = signInManager;
            _logger = logger;
        }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        public string ReturnUrl { get; set; }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        public IList<AuthenticationScheme> ExternalLogins { get; set; }

        public object ErrorMessage { get; set; }
        
        [BindProperty]
        [Required(ErrorMessage = "You have to enter a username")]
        [DataType(DataType.Text)]
        [Display(Name = "Username")]
        public string Username { get; set; }

        [BindProperty]
        [Required(ErrorMessage = "You have to enter a valid email address")]
        [EmailAddress]
        [Display(Name = "Email")]
        public string Email { get; set; }

        [BindProperty]
        [Required(ErrorMessage = "You have to enter a password")]
        [DataType(DataType.Password)]
        [Display(Name = "Password")]
        public string Password { get; set; }
        
        [BindProperty(Name = "password2")]
        [Required(ErrorMessage = "The two passwords do not match")]
        [DataType(DataType.Password)]
        [Display(Name = "Confirm Password")]
        [Compare("Password", ErrorMessage = "The two passwords do not match")]
        public string Password2 { get; set; }
        
        
        public async Task OnGetAsync(string returnUrl = null)
        {
            ReturnUrl = returnUrl;
            ExternalLogins = (
                await _signInManager.GetExternalAuthenticationSchemesAsync()
            ).ToList();
        }
        
        public async Task<IActionResult> OnPostAsync()
        {
            if (!ModelState.IsValid)
            {
                return Page();
            }

            // Check for existing username or email
            if (await _userManager.FindByNameAsync(Username) != null)
            {
                ModelState.AddModelError("Input.Username", "The username is already taken");
                return Page();
            }

            if (await _userManager.FindByEmailAsync(Email) != null)
            {
                ModelState.AddModelError("Input.Email", "The email is already registered");
                return Page();
            }

            // Create new user
            var newUser = new Author
            {
                UserName = Username,
                Email = Email,
                NormalizedEmail = Email.ToUpperInvariant(),
            };
            Console.WriteLine(newUser);
            var result = await _userManager.CreateAsync(newUser, Password);
            
            if (result.Succeeded)
            {
                TempData["FlashMessage"] = "You were successfully registered and can login now";
                return RedirectToPage("/Account/Login", new { area = "Identity" });
            }
            
            foreach (var error in result.Errors)
            {
                ModelState.AddModelError(string.Empty, error.Description);
            }

            return Page();
        }
        
        private Author CreateUser()
        {
            try
            {
                return Activator.CreateInstance<Author>();
            }
            catch
            {
                throw new InvalidOperationException(
                    $"Can't create an instance of 'Author'. "
                        + $"Ensure that 'Author' is not an abstract class and has a parameterless constructor, or alternatively "
                        + $"override the register page in /Areas/Identity/Pages/Account/Register.cshtml"
                );
            }
        }

        private IUserEmailStore<Author> GetEmailStore()
        {
            if (!_userManager.SupportsUserEmail)
            {
                throw new NotSupportedException(
                    "The default UI requires a user store with email support."
                );
            }
            return (IUserEmailStore<Author>)_userStore;
        }
    }
}
