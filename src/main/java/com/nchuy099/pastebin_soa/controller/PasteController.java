package com.nchuy099.pastebin_soa.controller;

import com.nchuy099.pastebin_soa.common.Visibility;
import com.nchuy099.pastebin_soa.dto.MonthlyStatsDTO;
import com.nchuy099.pastebin_soa.model.PasteEntity;
import com.nchuy099.pastebin_soa.service.PasteService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.time.Instant;
import java.util.Calendar;
import java.util.Date;
import java.util.List;

@Controller
@RequiredArgsConstructor
public class PasteController {

    private final PasteService pasteService;

    @GetMapping("/")
    public String showCreateForm(Model model) {
        model.addAttribute("error", null);
        return "index"; // Tương ứng với index.jsp
    }

    @PostMapping("/paste")
    public String createPaste(
            @RequestParam("content") String content,
            @RequestParam(value = "title", required = false) String title,
            @RequestParam(value = "language", required = false) String language,
            @RequestParam(value = "expires_in", required = false) Long expiresIn,
            @RequestParam(value = "visibility", required = false) String visibility,
            Model model) {
        try {
            if (content == null || content.trim().isEmpty()) {
                model.addAttribute("error", "Content is required");
                return "index";
            }

            PasteEntity paste = new PasteEntity();
            paste.setContent(content);
            paste.setTitle(title);
            paste.setLanguage(language);
            if (expiresIn != null && expiresIn > 0) {
                paste.setExpiresAt(Date.from(Instant.now().plusMillis(expiresIn * 60 * 1000)));
            }
            if (visibility != null) {
                paste.setVisibility(Visibility.valueOf(visibility.toUpperCase()));
            }

            PasteEntity result = pasteService.createPaste(paste);
            return "redirect:/paste/" + result.getId(); // Chuyển hướng đến trang xem paste
        } catch (Exception e) {
            System.err.println("Error creating paste: " + e.getMessage());
            model.addAttribute("error", "Failed to create paste");
            return "index";
        }
    }

    @GetMapping("/paste/{id}")
    public String getPaste(@PathVariable("id") String id, Model model) {
        try {
            PasteEntity paste = pasteService.getPasteById(id);
            model.addAttribute("paste", paste);
            return "paste"; // Tương ứng với paste.jsp
        } catch (Exception e) {
            System.err.println("Error retrieving paste: " + e.getMessage());
            String errorMessage = e.getMessage();
            if (errorMessage.contains("expired")) {
                model.addAttribute("error", errorMessage);
                model.addAttribute("pasteId", id);
                return "expired"; // Tương ứng với expired.jsp
            } else if (errorMessage.contains("not found")) {
                model.addAttribute("error", errorMessage);
                model.addAttribute("pasteId", id);
                return "index"; // 404 -> trả về index với lỗi
            } else {
                model.addAttribute("error", "Failed to retrieve paste");
                model.addAttribute("pasteId", id);
                return "index"; // 500 -> trả về index với lỗi
            }
        }
    }

    @GetMapping("/public")
    public String getPublicPastes(Model model) {
        try {
            List<PasteEntity> pastes = pasteService.getPublicPastes();
            model.addAttribute("pastes", pastes);
            return "public"; // Tương ứng với public.jsp
        } catch (Exception e) {
            System.err.println("Error fetching public pastes: " + e.getMessage());
            model.addAttribute("pastes", List.of()); // Danh sách rỗng nếu lỗi
            model.addAttribute("error", "Failed to fetch public pastes");
            return "public";
        }
    }

    @GetMapping(value = { "/stats", "/stats/{month}" })
    public String getMonthlyStats(@PathVariable(value = "month", required = false) String month, Model model) {
        try {
            Date targetMonth;
            if (month == null) {
                Calendar cal = Calendar.getInstance();
                cal.set(Calendar.DAY_OF_MONTH, 1);
                cal.set(Calendar.HOUR_OF_DAY, 0);
                cal.set(Calendar.MINUTE, 0);
                cal.set(Calendar.SECOND, 0);
                cal.set(Calendar.MILLISECOND, 0);
                targetMonth = cal.getTime();
            } else {
                if (!month.matches("\\d{4}-\\d{2}")) {
                    throw new Exception("Invalid month format. Use YYYY-MM");
                }
                SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM");
                targetMonth = sdf.parse(month);
            }

            MonthlyStatsDTO stats = pasteService.getMonthlyStats(targetMonth);
            model.addAttribute("stats", stats);
            model.addAttribute("error", null);
            return "stats";
        } catch (ParseException e) {
            System.err.println("Error parsing month: " + e.getMessage());
            model.addAttribute("stats", null);
            model.addAttribute("error", "Invalid month format");
            return "stats";
        } catch (Exception e) {
            System.err.println("Error fetching monthly stats: " + e.getMessage());
            model.addAttribute("stats", null);
            model.addAttribute("error", e.getMessage() != null ? e.getMessage() : "Failed to fetch statistics");
            return "stats";
        }
    }
}